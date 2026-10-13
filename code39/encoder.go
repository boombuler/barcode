// Package code39 can create Code39 barcodes
package code39

import (
	"errors"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type encodeInfo struct {
	value int
	data  []bool
}

var encodeTable = map[rune]encodeInfo{
	'0': {0, []bool{true, false, true, false, false, true, true, false, true, true, false, true}},
	'1': {1, []bool{true, true, false, true, false, false, true, false, true, false, true, true}},
	'2': {2, []bool{true, false, true, true, false, false, true, false, true, false, true, true}},
	'3': {3, []bool{true, true, false, true, true, false, false, true, false, true, false, true}},
	'4': {4, []bool{true, false, true, false, false, true, true, false, true, false, true, true}},
	'5': {5, []bool{true, true, false, true, false, false, true, true, false, true, false, true}},
	'6': {6, []bool{true, false, true, true, false, false, true, true, false, true, false, true}},
	'7': {7, []bool{true, false, true, false, false, true, false, true, true, false, true, true}},
	'8': {8, []bool{true, true, false, true, false, false, true, false, true, true, false, true}},
	'9': {9, []bool{true, false, true, true, false, false, true, false, true, true, false, true}},
	'A': {10, []bool{true, true, false, true, false, true, false, false, true, false, true, true}},
	'B': {11, []bool{true, false, true, true, false, true, false, false, true, false, true, true}},
	'C': {12, []bool{true, true, false, true, true, false, true, false, false, true, false, true}},
	'D': {13, []bool{true, false, true, false, true, true, false, false, true, false, true, true}},
	'E': {14, []bool{true, true, false, true, false, true, true, false, false, true, false, true}},
	'F': {15, []bool{true, false, true, true, false, true, true, false, false, true, false, true}},
	'G': {16, []bool{true, false, true, false, true, false, false, true, true, false, true, true}},
	'H': {17, []bool{true, true, false, true, false, true, false, false, true, true, false, true}},
	'I': {18, []bool{true, false, true, true, false, true, false, false, true, true, false, true}},
	'J': {19, []bool{true, false, true, false, true, true, false, false, true, true, false, true}},
	'K': {20, []bool{true, true, false, true, false, true, false, true, false, false, true, true}},
	'L': {21, []bool{true, false, true, true, false, true, false, true, false, false, true, true}},
	'M': {22, []bool{true, true, false, true, true, false, true, false, true, false, false, true}},
	'N': {23, []bool{true, false, true, false, true, true, false, true, false, false, true, true}},
	'O': {24, []bool{true, true, false, true, false, true, true, false, true, false, false, true}},
	'P': {25, []bool{true, false, true, true, false, true, true, false, true, false, false, true}},
	'Q': {26, []bool{true, false, true, false, true, false, true, true, false, false, true, true}},
	'R': {27, []bool{true, true, false, true, false, true, false, true, true, false, false, true}},
	'S': {28, []bool{true, false, true, true, false, true, false, true, true, false, false, true}},
	'T': {29, []bool{true, false, true, false, true, true, false, true, true, false, false, true}},
	'U': {30, []bool{true, true, false, false, true, false, true, false, true, false, true, true}},
	'V': {31, []bool{true, false, false, true, true, false, true, false, true, false, true, true}},
	'W': {32, []bool{true, true, false, false, true, true, false, true, false, true, false, true}},
	'X': {33, []bool{true, false, false, true, false, true, true, false, true, false, true, true}},
	'Y': {34, []bool{true, true, false, false, true, false, true, true, false, true, false, true}},
	'Z': {35, []bool{true, false, false, true, true, false, true, true, false, true, false, true}},
	'-': {36, []bool{true, false, false, true, false, true, false, true, true, false, true, true}},
	'.': {37, []bool{true, true, false, false, true, false, true, false, true, true, false, true}},
	' ': {38, []bool{true, false, false, true, true, false, true, false, true, true, false, true}},
	'$': {39, []bool{true, false, false, true, false, false, true, false, false, true, false, true}},
	'/': {40, []bool{true, false, false, true, false, false, true, false, true, false, false, true}},
	'+': {41, []bool{true, false, false, true, false, true, false, false, true, false, false, true}},
	'%': {42, []bool{true, false, true, false, false, true, false, false, true, false, false, true}},
	'*': {-1, []bool{true, false, false, true, false, true, true, false, true, true, false, true}},
}

var extendedTable = map[rune]string{
	0: `%U`, 1: `$A`, 2: `$B`, 3: `$C`, 4: `$D`, 5: `$E`, 6: `$F`, 7: `$G`, 8: `$H`, 9: `$I`, 10: `$J`,
	11: `$K`, 12: `$L`, 13: `$M`, 14: `$N`, 15: `$O`, 16: `$P`, 17: `$Q`, 18: `$R`, 19: `$S`, 20: `$T`,
	21: `$U`, 22: `$V`, 23: `$W`, 24: `$X`, 25: `$Y`, 26: `$Z`, 27: `%A`, 28: `%B`, 29: `%C`, 30: `%D`,
	31: `%E`, 33: `/A`, 34: `/B`, 35: `/C`, 36: `/D`, 37: `/E`, 38: `/F`, 39: `/G`, 40: `/H`, 41: `/I`,
	42: `/J`, 43: `/K`, 44: `/L`, 47: `/O`, 58: `/Z`, 59: `%F`, 60: `%G`, 61: `%H`, 62: `%I`, 63: `%J`,
	64: `%V`, 91: `%K`, 92: `%L`, 93: `%M`, 94: `%N`, 95: `%O`, 96: `%W`, 97: `+A`, 98: `+B`, 99: `+C`,
	100: `+D`, 101: `+E`, 102: `+F`, 103: `+G`, 104: `+H`, 105: `+I`, 106: `+J`, 107: `+K`, 108: `+L`,
	109: `+M`, 110: `+N`, 111: `+O`, 112: `+P`, 113: `+Q`, 114: `+R`, 115: `+S`, 116: `+T`, 117: `+U`,
	118: `+V`, 119: `+W`, 120: `+X`, 121: `+Y`, 122: `+Z`, 123: `%P`, 124: `%Q`, 125: `%R`, 126: `%S`,
	127: `%T`,
}

func getChecksum(content string) rune {
	sum := 0
	for _, r := range content {
		info, ok := encodeTable[r]
		if !ok || info.value < 0 {
			return 0
		}

		sum += info.value
	}

	sum = sum % 43
	for r, v := range encodeTable {
		if v.value == sum {
			return r
		}
	}
	return 0
}

func prepare(content string) (string, error) {
	result := ""
	for _, r := range content {
		if r > 127 {
			return "", errors.New("only ASCII strings can be encoded")
		}
		val, ok := extendedTable[r]
		if ok {
			result += val
		} else {
			result += string([]rune{r})
		}
	}
	return result, nil
}

// Encode returns a code39 barcode for the given content and color scheme
// if includeChecksum is set to true, a checksum character is calculated and added to the content
func EncodeWithColor(content string, includeChecksum bool, fullASCIIMode bool, color barcode.ColorScheme) (barcode.BarcodeIntCS, error) {
	if fullASCIIMode {
		var err error
		content, err = prepare(content)
		if err != nil {
			return nil, err
		}
	} else if strings.ContainsRune(content, '*') {
		return nil, errors.New("invalid data! try full ascii mode")
	}

	checkChar := getChecksum(content)
	data := "*" + content
	if includeChecksum {
		data += string(checkChar)
	}
	data += "*"

	result := new(utils.BitList)

	for i, r := range data {
		if i != 0 {
			result.AddBit(false)
		}

		info, ok := encodeTable[r]
		if !ok {
			return nil, errors.New("invalid data! try full ascii mode")
		}
		result.AddBit(info.data...)
	}

	checkSum := encodeTable[checkChar].value
	return utils.New1DCodeIntCheckSumWithColor(barcode.TypeCode39, content, result, checkSum, color), nil
}

// Encode returns a code39 barcode for the given content
// if includeChecksum is set to true, a checksum character is calculated and added to the content
func Encode(content string, includeChecksum bool, fullASCIIMode bool) (barcode.BarcodeIntCS, error) {
	return EncodeWithColor(content, includeChecksum, fullASCIIMode, barcode.ColorScheme16)
}
