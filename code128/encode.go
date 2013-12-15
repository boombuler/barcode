package code128

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
	"unicode/utf8"
)

func strToRunes(str string) []rune {
	result := make([]rune, utf8.RuneCountInString(str))
	i := 0
	for _, r := range str {
		result[i] = r
		i++
	}
	return result
}

func Encode(content string) (barcode.Barcode, error) {
	contentRunes := strToRunes(content)
	idxList := encodeBTable(contentRunes)

	if idxList == nil {
		return nil, fmt.Errorf("\"%s\" could not be encoded", content)
	}

	result := &code{new(utils.BitList), content}
	sum := 0
	for i, idx := range idxList.GetBytes() {
		if i == 0 {
			sum = int(idx)
		} else {
			sum += i * int(idx)
		}
		result.AddRange(encodingTable[idx])
	}
	result.AddRange(encodingTable[sum%103])
	result.AddRange(encodingTable[stopSymbol])
	return result, nil
}
