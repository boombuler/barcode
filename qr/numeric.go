package qr

import (
	"errors"
	"fmt"

	"github.com/boombuler/barcode/utils"
)

func encodeNumeric(content string, ecl ErrorCorrectionLevel) (*utils.BitList, *versionInfo, error) {
	if !isNumeric(content) {
		return nil, nil, fmt.Errorf("%q can not be encoded as %s", content, Numeric)
	}

	contentBitCount := (len(content) / 3) * 10
	switch len(content) % 3 {
	case 1:
		contentBitCount += 4
	case 2:
		contentBitCount += 7
	}
	vi := findSmallestVersionInfo(ecl, numericMode, contentBitCount)
	if vi == nil {
		return nil, nil, errors.New("too much data to encode")
	}
	res := new(utils.BitList)
	res.AddBits(int(numericMode), 4)
	res.AddBits(len(content), vi.charCountBits(numericMode))

	for pos := 0; pos < len(content); pos += 3 {
		var group string
		if pos+3 <= len(content) {
			group = content[pos : pos+3]
		} else {
			group = content[pos:]
		}

		var value int
		for _, c := range group {
			value = value*10 + int(c-'0')
		}

		var bitCnt byte
		switch len(group) % 3 {
		case 0:
			bitCnt = 10
		case 1:
			bitCnt = 4
		case 2:
			bitCnt = 7
		}

		res.AddBits(value, bitCnt)
	}

	addPaddingAndTerminator(res, vi)
	return res, vi, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
