package code39

import (
	"errors"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
	"strings"
)

type encodeInfo struct {
	value int
	data  []bool
}

var encodeTable map[rune]encodeInfo = map[rune]encodeInfo{
	'0': encodeInfo{0, []bool{true, false, true, false, false, true, true, false, true, true, false, true}},
	'1': encodeInfo{1, []bool{true, true, false, true, false, false, true, false, true, false, true, true}},
	'2': encodeInfo{2, []bool{true, false, true, true, false, false, true, false, true, false, true, true}},
	'3': encodeInfo{3, []bool{true, true, false, true, true, false, false, true, false, true, false, true}},
	'4': encodeInfo{4, []bool{true, false, true, false, false, true, true, false, true, false, true, true}},
	'5': encodeInfo{5, []bool{true, true, false, true, false, false, true, true, false, true, false, true}},
	'6': encodeInfo{6, []bool{true, false, true, true, false, false, true, true, false, true, false, true}},
	'7': encodeInfo{7, []bool{true, false, true, false, false, true, false, true, true, false, true, true}},
	'8': encodeInfo{8, []bool{true, true, false, true, false, false, true, false, true, true, false, true}},
	'9': encodeInfo{9, []bool{true, false, true, true, false, false, true, false, true, true, false, true}},
	'A': encodeInfo{10, []bool{true, true, false, true, false, true, false, false, true, false, true, true}},
	'B': encodeInfo{11, []bool{true, false, true, true, false, true, false, false, true, false, true, true}},
	'C': encodeInfo{12, []bool{true, true, false, true, true, false, true, false, false, true, false, true}},
	'D': encodeInfo{13, []bool{true, false, true, false, true, true, false, false, true, false, true, true}},
	'E': encodeInfo{14, []bool{true, true, false, true, false, true, true, false, false, true, false, true}},
	'F': encodeInfo{15, []bool{true, false, true, true, false, true, true, false, false, true, false, true}},
	'G': encodeInfo{16, []bool{true, false, true, false, true, false, false, true, true, false, true, true}},
	'H': encodeInfo{17, []bool{true, true, false, true, false, true, false, false, true, true, false, true}},
	'I': encodeInfo{18, []bool{true, false, true, true, false, true, false, false, true, true, false, true}},
	'J': encodeInfo{19, []bool{true, false, true, false, true, true, false, false, true, true, false, true}},
	'K': encodeInfo{20, []bool{true, true, false, true, false, true, false, true, false, false, true, true}},
	'L': encodeInfo{21, []bool{true, false, true, true, false, true, false, true, false, false, true, true}},
	'M': encodeInfo{22, []bool{true, true, false, true, true, false, true, false, true, false, false, true}},
	'N': encodeInfo{23, []bool{true, false, true, false, true, true, false, true, false, false, true, true}},
	'O': encodeInfo{24, []bool{true, true, false, true, false, true, true, false, true, false, false, true}},
	'P': encodeInfo{25, []bool{true, false, true, true, false, true, true, false, true, false, false, true}},
	'Q': encodeInfo{26, []bool{true, false, true, false, true, false, true, true, false, false, true, true}},
	'R': encodeInfo{27, []bool{true, true, false, true, false, true, false, true, true, false, false, true}},
	'S': encodeInfo{28, []bool{true, false, true, true, false, true, false, true, true, false, false, true}},
	'T': encodeInfo{29, []bool{true, false, true, false, true, true, false, true, true, false, false, true}},
	'U': encodeInfo{30, []bool{true, true, false, false, true, false, true, false, true, false, true, true}},
	'V': encodeInfo{31, []bool{true, false, false, true, true, false, true, false, true, false, true, true}},
	'W': encodeInfo{32, []bool{true, true, false, false, true, true, false, true, false, true, false, true}},
	'X': encodeInfo{33, []bool{true, false, false, true, false, true, true, false, true, false, true, true}},
	'Y': encodeInfo{34, []bool{true, true, false, false, true, false, true, true, false, true, false, true}},
	'Z': encodeInfo{35, []bool{true, false, false, true, true, false, true, true, false, true, false, true}},
	'-': encodeInfo{36, []bool{true, false, false, true, false, true, false, true, true, false, true, true}},
	'.': encodeInfo{37, []bool{true, true, false, false, true, false, true, false, true, true, false, true}},
	' ': encodeInfo{38, []bool{true, false, false, true, true, false, true, false, true, true, false, true}},
	'$': encodeInfo{39, []bool{true, false, false, true, false, false, true, false, false, true, false, true}},
	'/': encodeInfo{40, []bool{true, false, false, true, false, false, true, false, true, false, false, true}},
	'+': encodeInfo{41, []bool{true, false, false, true, false, true, false, false, true, false, false, true}},
	'%': encodeInfo{42, []bool{true, false, true, false, false, true, false, false, true, false, false, true}},
	'*': encodeInfo{-1, []bool{true, false, false, true, false, true, true, false, true, true, false, true}},
}

func getChecksum(content string) string {
	sum := 0
	for _, r := range content {
		info, ok := encodeTable[r]
		if !ok || info.value < 0 {
			return "#"
		}

		sum += info.value
	}

	sum = sum % 43
	for r, v := range encodeTable {
		if v.value == sum {
			return string(r)
		}
	}
	return "#"
}

// encodes the given string as a code39 barcode
// if includeChecksum is set to true, a checksum character is calculated and added to the content
func Encode(content string, includeChecksum bool) (barcode.Barcode, error) {
	if strings.ContainsRune(content, '*') {
		return nil, errors.New("invalid data")
	}

	data := "*" + content
	if includeChecksum {
		data += getChecksum(content)
	}
	data += "*"

	result := new(utils.BitList)

	for i, r := range data {
		if i != 0 {
			result.AddBit(false)
		}

		info, ok := encodeTable[r]
		if !ok {
			return nil, errors.New("invalid data")
		}
		result.AddBit(info.data...)
	}

	return utils.New1DCode("Code 39", content, result), nil
}
