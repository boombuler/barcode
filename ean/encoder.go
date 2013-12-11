package ean

import (
	"errors"
	"github.com/boombuler/barcode"
)

type encodedNumber struct {
	LeftOdd  []bool
	LeftEven []bool
	Right    []bool
	CheckSum []bool
}

var encoderTable map[rune]encodedNumber = map[rune]encodedNumber{
	'0': encodedNumber{
		[]bool{false, false, false, true, true, false, true},
		[]bool{false, true, false, false, true, true, true},
		[]bool{true, true, true, false, false, true, false},
		[]bool{false, false, false, false, false, false},
	},
	'1': encodedNumber{
		[]bool{false, false, true, true, false, false, true},
		[]bool{false, true, true, false, false, true, true},
		[]bool{true, true, false, false, true, true, false},
		[]bool{false, false, true, false, true, true},
	},
	'2': encodedNumber{
		[]bool{false, false, true, false, false, true, true},
		[]bool{false, false, true, true, false, true, true},
		[]bool{true, true, false, true, true, false, false},
		[]bool{false, false, true, true, false, true},
	},
	'3': encodedNumber{
		[]bool{false, true, true, true, true, false, true},
		[]bool{false, true, false, false, false, false, true},
		[]bool{true, false, false, false, false, true, false},
		[]bool{false, false, true, true, true, false},
	},
	'4': encodedNumber{
		[]bool{false, true, false, false, false, true, true},
		[]bool{false, false, true, true, true, false, true},
		[]bool{true, false, true, true, true, false, false},
		[]bool{false, true, false, false, true, true},
	},
	'5': encodedNumber{
		[]bool{false, true, true, false, false, false, true},
		[]bool{false, true, true, true, false, false, true},
		[]bool{true, false, false, true, true, true, false},
		[]bool{false, true, true, false, false, true},
	},
	'6': encodedNumber{
		[]bool{false, true, false, true, true, true, true},
		[]bool{false, false, false, false, true, false, true},
		[]bool{true, false, true, false, false, false, false},
		[]bool{false, true, true, true, false, false},
	},
	'7': encodedNumber{
		[]bool{false, true, true, true, false, true, true},
		[]bool{false, false, true, false, false, false, true},
		[]bool{true, false, false, false, true, false, false},
		[]bool{false, true, false, true, false, true},
	},
	'8': encodedNumber{
		[]bool{false, true, true, false, true, true, true},
		[]bool{false, false, false, true, false, false, true},
		[]bool{true, false, false, true, false, false, false},
		[]bool{false, true, false, true, true, false},
	},
	'9': encodedNumber{
		[]bool{false, false, false, true, false, true, true},
		[]bool{false, false, true, false, true, true, true},
		[]bool{true, true, true, false, true, false, false},
		[]bool{false, true, true, false, true, false},
	},
}

func runeToInt(r rune) int {
	switch r {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	}
	return -1
}

func intToRune(i int) rune {
	switch i {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	case 6:
		return '6'
	case 7:
		return '7'
	case 8:
		return '8'
	case 9:
		return '9'
	}
	return 'F'
}

func calcCheckNum(code string) rune {
	x3 := len(code) == 7
	sum := 0
	for _, r := range code {
		curNum := runeToInt(r)
		if curNum < 0 || curNum > 9 {
			return 'B'
		}
		if x3 {
			curNum = curNum * 3
		}
		x3 = !x3
		sum += curNum
	}

	return intToRune((10 - (sum % 10)) % 10)
}

func encodeEAN8(code string, result *eancode) bool {
	pos := 0
	appendBit := func(b bool) {
		result.SetBit(pos, b)
		pos++
	}

	appendBit(true)
	appendBit(false)
	appendBit(true)

	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return false
		}
		var data []bool
		if cpos < 4 {
			data = num.LeftOdd
		} else {
			data = num.Right
		}

		if cpos == 4 {
			appendBit(false)
			appendBit(true)
			appendBit(false)
			appendBit(true)
			appendBit(false)
		}
		for _, bit := range data {
			appendBit(bit)
		}
	}

	appendBit(true)
	appendBit(false)
	appendBit(true)
	return true
}

func encodeEAN13(code string, result *eancode) bool {
	pos := 0
	appendBit := func(b bool) {
		result.SetBit(pos, b)
		pos++
	}

	appendBit(true)
	appendBit(false)
	appendBit(true)

	var firstNum []bool
	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return false
		}
		if cpos == 0 {
			firstNum = num.CheckSum
			continue
		}

		var data []bool
		if cpos < 7 { // Left
			if firstNum[cpos-1] {
				data = num.LeftEven
			} else {
				data = num.LeftOdd
			}
		} else {
			data = num.Right
		}

		if cpos == 7 {
			appendBit(false)
			appendBit(true)
			appendBit(false)
			appendBit(true)
			appendBit(false)
		}

		for _, bit := range data {
			appendBit(bit)
		}
	}
	appendBit(true)
	appendBit(false)
	appendBit(true)
	return true
}

// encodes the given EAN 8 or EAN 13 number to a barcode image
func Encode(code string) (barcode.Barcode, error) {
	if len(code) == 7 || len(code) == 12 {
		code += string(calcCheckNum(code))
	} else if len(code) == 8 || len(code) == 13 {
		check := code[0 : len(code)-1]
		check += string(calcCheckNum(check))
		if check != code {
			return nil, errors.New("checksum missmatch!")
		}
	}
	ean8 := false
	if len(code) == 8 {
		ean8 = true

	} else if len(code) != 13 {
		return nil, errors.New("invalid ean code data")
	}
	result := newEANCode(ean8)
	if (ean8 && encodeEAN8(code, result)) || (!ean8 && encodeEAN13(code, result)) {
		result.content = code
		return result, nil
	}

	return nil, errors.New("ean code contains invalid characters")
}
