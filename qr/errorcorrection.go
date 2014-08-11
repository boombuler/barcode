package qr

import (
	"github.com/boombuler/barcode/utils"
)

type errorCorrection struct {
	fld       *utils.GaloisField
	polynomes map[byte][]byte
}

var ec = newGF()

func newGF() *errorCorrection {
	return &errorCorrection{utils.NewGaloisField(285), make(map[byte][]byte)}
}

func (ec *errorCorrection) getPolynomial(eccc byte) []byte {
	_, ok := ec.polynomes[eccc]
	if !ok {
		if eccc == 1 {
			ec.polynomes[eccc] = []byte{0, 0}
		} else {
			b1 := ec.getPolynomial(eccc - 1)
			result := make([]byte, eccc+1)
			for x := 0; x < len(b1); x++ {
				tmp1 := (int(b1[x]) + int(eccc-1)) % 255
				if x == 0 {
					result[x] = b1[x]
				} else {
					tmp0 := int(ec.fld.ALogTbl[result[x]]) ^ int(ec.fld.ALogTbl[b1[x]])
					result[x] = byte(ec.fld.LogTbl[tmp0])
				}
				result[x+1] = byte(tmp1)
			}
			ec.polynomes[eccc] = result

		}
	}
	return ec.polynomes[eccc]
}

func (ec *errorCorrection) calcECC(data []byte, eccCount byte) []byte {
	tmp := make([]byte, len(data)+int(eccCount))
	copy(tmp, data)
	generator := ec.getPolynomial(eccCount)

	for i := 0; i < len(data); i++ {
		alpha := ec.fld.LogTbl[tmp[i]]
		for j := 0; j < len(generator); j++ {
			idx := (int(alpha) + int(generator[j])) % 255
			polyJ := ec.fld.ALogTbl[idx]
			tmp[i+j] = byte(ec.fld.AddOrSub(int(tmp[i+j]), polyJ))
		}
	}

	return tmp[len(data):]
}
