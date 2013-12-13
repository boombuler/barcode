package utils

type GaloisField struct {
	ALogTbl []int
	LogTbl  []int
}

func NewGaloisField(pp int) *GaloisField {
	result := new(GaloisField)
	fldSize := 256

	result.ALogTbl = make([]int, fldSize)
	result.LogTbl = make([]int, fldSize)

	x := 1
	for i := 0; i < fldSize; i++ {
		result.ALogTbl[i] = x
		x = x * 2
		if x >= fldSize {
			x = (x ^ pp) & (fldSize - 1)
		}
	}

	for i := 0; i < fldSize; i++ {
		result.LogTbl[result.ALogTbl[i]] = int(i)
	}

	return result
}

func (gf *GaloisField) AddOrSub(a, b int) int {
	return a ^ b
}

func (gf *GaloisField) Multiply(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return gf.ALogTbl[(gf.LogTbl[a]+gf.LogTbl[b])%255]
}

func (gf *GaloisField) Divide(a, b int) int {
	if b == 0 {
		panic("divide by zero")
	} else if a == 0 {
		return 0
	}
	return gf.ALogTbl[(gf.LogTbl[a]-gf.LogTbl[b])%255]
}
