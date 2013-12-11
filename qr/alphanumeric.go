package qr

import (
	"errors"
	"fmt"
	"github.com/boombuler/barcode"
)

var alphaNumericTable map[byte]int = map[byte]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19,
	'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29,
	'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35, ' ': 36, '$': 37, '%': 38, '*': 39,
	'+': 40, '-': 41, '.': 42, '/': 43, ':': 44,
}

func encodeAlphaNumeric(content string, ecl ErrorCorrectionLevel) (*barcode.BitList, *versionInfo, error) {

	contentLenIsOdd := len(content)%2 == 1
	contentBitCount := (len(content) / 2) * 11
	if contentLenIsOdd {
		contentBitCount += 6
	}
	vi := findSmallestVersionInfo(ecl, alphaNumericMode, contentBitCount)
	if vi == nil {
		return nil, nil, errors.New("To much data to encode")
	}

	res := new(barcode.BitList)
	res.AddBits(int(alphaNumericMode), 4)
	res.AddBits(len(content), vi.charCountBits(alphaNumericMode))

	for idx := 0; idx < len(content)/2; idx++ {
		c1, ok1 := alphaNumericTable[content[idx*2]]
		c2, ok2 := alphaNumericTable[content[(idx*2)+1]]
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("\"%s\" can not be encoded as %s", content, AlphaNumeric)
		}
		res.AddBits(c1*45+c2, 11)
	}
	if contentLenIsOdd {
		c1, ok := alphaNumericTable[content[len(content)-1]]
		if !ok {
			return nil, nil, fmt.Errorf("\"%s\" can not be encoded as %s", content, AlphaNumeric)
		}
		res.AddBits(c1, 6)
	}

	addPaddingAndTerminator(res, vi)

	return res, vi, nil
}
