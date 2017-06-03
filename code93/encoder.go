// Package code39 can create Code39 barcodes
package code93

import (
	"errors"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type encodeInfo struct {
	value int
	data  int
}

const (
	// Special Function 1 ($)
	FNC1 = '\u00f1'
	// Special Function 2 (%)
	FNC2 = '\u00f2'
	// Special Function 3 (/)
	FNC3 = '\u00f3'
	// Special Function 4 (+)
	FNC4 = '\u00f4'
)

var encodeTable = map[rune]encodeInfo{
	'0': encodeInfo{0, 0x114}, '1': encodeInfo{1, 0x148}, '2': encodeInfo{2, 0x144},
	'3': encodeInfo{3, 0x142}, '4': encodeInfo{4, 0x128}, '5': encodeInfo{5, 0x124},
	'6': encodeInfo{6, 0x122}, '7': encodeInfo{7, 0x150}, '8': encodeInfo{8, 0x112},
	'9': encodeInfo{9, 0x10A}, 'A': encodeInfo{10, 0x1A8}, 'B': encodeInfo{11, 0x1A4},
	'C': encodeInfo{12, 0x1A2}, 'D': encodeInfo{13, 0x194}, 'E': encodeInfo{14, 0x192},
	'F': encodeInfo{15, 0x18A}, 'G': encodeInfo{16, 0x168}, 'H': encodeInfo{17, 0x164},
	'I': encodeInfo{18, 0x162}, 'J': encodeInfo{19, 0x134}, 'K': encodeInfo{20, 0x11A},
	'L': encodeInfo{21, 0x158}, 'M': encodeInfo{22, 0x14C}, 'N': encodeInfo{23, 0x146},
	'O': encodeInfo{24, 0x12C}, 'P': encodeInfo{25, 0x116}, 'Q': encodeInfo{26, 0x1B4},
	'R': encodeInfo{27, 0x1B2}, 'S': encodeInfo{28, 0x1AC}, 'T': encodeInfo{29, 0x1A6},
	'U': encodeInfo{30, 0x196}, 'V': encodeInfo{31, 0x19A}, 'W': encodeInfo{32, 0x16C},
	'X': encodeInfo{33, 0x166}, 'Y': encodeInfo{34, 0x136}, 'Z': encodeInfo{35, 0x13A},
	'-': encodeInfo{36, 0x12E}, '.': encodeInfo{37, 0x1D4}, ' ': encodeInfo{38, 0x1D2},
	'$': encodeInfo{39, 0x1CA}, '/': encodeInfo{40, 0x16E}, '+': encodeInfo{41, 0x176},
	'%': encodeInfo{42, 0x1AE}, FNC1: encodeInfo{43, 0x126}, FNC2: encodeInfo{44, 0x1DA},
	FNC3: encodeInfo{45, 0x1D6}, FNC4: encodeInfo{46, 0x132}, '*': encodeInfo{47, 0x15E},
}

// Encode returns a code93 barcode for the given content
func Encode(content string) (barcode.Barcode, error) {
	if strings.ContainsRune(content, '*') {
		return nil, errors.New("invalid data! content may not contain '*'")
	}

	data := content + string(getChecksum(content, 20))
	data += string(getChecksum(data, 15))

	data = "*" + data + "*"
	result := new(utils.BitList)

	for _, r := range data {
		info, ok := encodeTable[r]
		if !ok {
			return nil, errors.New("invalid data!")
		}
		result.AddBits(info.data, 9)
	}
	result.AddBit(true)

	return utils.New1DCode("Code 93", content, result), nil
}

func reverse(value string) string {
	data := []rune(value)
	result := []rune{}
	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}
	return string(result)
}

func getChecksum(content string, maxWeight int) rune {
	weight := 1
	total := 0

	for _, r := range reverse(content) {
		info, ok := encodeTable[r]
		if !ok {
			return ' '
		}
		total += info.value * weight
		if weight++; weight > maxWeight {
			weight = 1
		}
	}
	total = total % 47
	for r, info := range encodeTable {
		if info.value == total {
			return r
		}
	}
	return ' '
}
