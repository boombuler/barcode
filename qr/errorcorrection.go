package qr

type galoisField struct {
	aLogTbl   []byte
	logTbl    []byte
	polynomes map[byte][]byte
}

var gf *galoisField = newGF()

func newGF() *galoisField {
	result := new(galoisField)
	result.polynomes = make(map[byte][]byte)
	result.aLogTbl = make([]byte, 255)
	result.logTbl = make([]byte, 256)

	result.aLogTbl[0] = 1

	x := 1
	for i := 1; i < 255; i++ {
		x = x * 2
		if x > 255 {
			x = x ^ 285
		}
		result.aLogTbl[i] = byte(x)
	}

	for i := 1; i < 255; i++ {
		result.logTbl[result.aLogTbl[i]] = byte(i)
	}

	return result
}

func (gf *galoisField) getPolynom(eccc byte) []byte {
	_, ok := gf.polynomes[eccc]
	if !ok {
		if eccc == 1 {
			gf.polynomes[eccc] = []byte{0, 0}
		} else {
			b1 := gf.getPolynom(eccc - 1)
			result := make([]byte, eccc+1)
			for x := 0; x < len(b1); x++ {
				tmp1 := (int(b1[x]) + int(eccc-1)) % 255
				if x == 0 {
					result[x] = b1[x]
				} else {
					tmp0 := int(gf.aLogTbl[result[x]]) ^ int(gf.aLogTbl[b1[x]])
					result[x] = gf.logTbl[tmp0]
				}
				result[x+1] = byte(tmp1)
			}
			gf.polynomes[eccc] = result

		}
	}
	return gf.polynomes[eccc]
}

func (gf *galoisField) calcECC(data []byte, eccCount byte) []byte {
	tmp := make([]byte, len(data)+int(eccCount))
	copy(tmp, data)
	generator := gf.getPolynom(eccCount)

	for i := 0; i < len(data); i++ {
		alpha := gf.logTbl[tmp[i]]
		for j := 0; j < len(generator); j++ {
			idx := (int(alpha) + int(generator[j])) % 255
			polyJ := gf.aLogTbl[idx]
			tmp[i+j] = (tmp[i+j] ^ polyJ)
		}
	}

	return tmp[len(data):]
}
