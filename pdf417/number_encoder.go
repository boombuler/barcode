package pdf417

import (
	"errors"
	"math/big"

	"github.com/boombuler/barcode/utils"
)

type numberEncoder struct{}

func newNumberEncoder() *numberEncoder {
	return new(numberEncoder)
}

func (encoder numberEncoder) CanEncode(char rune) bool {
	return utils.RuneToInt(char) != -1
}

func (encoder numberEncoder) GetSwitchCode(data string) int {
	return switchCodeNumber
}

func (encoder numberEncoder) Encode(digits string) ([]int, error) {
	digitCount := len(digits)
	chunkCount := digitCount / 44
	if digitCount%44 != 0 {
		chunkCount++
	}

	codeWords := []int{}

	for i := 0; i < chunkCount; i++ {
		start := i * 44
		end := start + 44
		if end > digitCount {
			end = digitCount
		}
		chunk := digits[start:end]

		cws, err := encodeChunk(chunk)
		if err != nil {
			return codeWords, err
		}

		codeWords = append(codeWords, cws...)
	}

	return codeWords, nil
}

func encodeChunk(chunkInput string) ([]int, error) {
	chunk := big.NewInt(0)

	_, ok := chunk.SetString("1"+chunkInput, 10)

	if !ok {
		return nil, errors.New("Failed converting")
	}

	cws := []int{}

	for chunk.Cmp(big.NewInt(0)) > 0 {
		newChunk, cw := chunk.DivMod(
			chunk,
			big.NewInt(900),
			big.NewInt(0),
		)

		chunk = newChunk

		cws = append([]int{int(cw.Int64())}, cws...)
	}

	return cws, nil
}
