// Package codabar can create Codabar barcodes
package codabar

import (
	"fmt"
	"regexp"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

var encodingTable = map[rune][]bool{
	'0': {true, false, true, false, true, false, false, true, true},
	'1': {true, false, true, false, true, true, false, false, true},
	'2': {true, false, true, false, false, true, false, true, true},
	'3': {true, true, false, false, true, false, true, false, true},
	'4': {true, false, true, true, false, true, false, false, true},
	'5': {true, true, false, true, false, true, false, false, true},
	'6': {true, false, false, true, false, true, false, true, true},
	'7': {true, false, false, true, false, true, true, false, true},
	'8': {true, false, false, true, true, false, true, false, true},
	'9': {true, true, false, true, false, false, true, false, true},
	'-': {true, false, true, false, false, true, true, false, true},
	'$': {true, false, true, true, false, false, true, false, true},
	':': {true, true, false, true, false, true, true, false, true, true},
	'/': {true, true, false, true, true, false, true, false, true, true},
	'.': {true, true, false, true, true, false, true, true, false, true},
	'+': {true, false, true, true, false, true, true, false, true, true},
	'A': {true, false, true, true, false, false, true, false, false, true},
	'B': {true, false, false, true, false, false, true, false, true, true},
	'C': {true, false, true, false, false, true, false, false, true, true},
	'D': {true, false, true, false, false, true, true, false, false, true},
}

// Encode creates a codabar barcode for the given content and color scheme
func EncodeWithColor(content string, color barcode.ColorScheme) (barcode.Barcode, error) {
	checkValid, _ := regexp.Compile(`[ABCD][0123456789\-\$\:/\.\+]*[ABCD]$`)
	if content == "!" || checkValid.ReplaceAllString(content, "!") != "!" {
		return nil, fmt.Errorf("can not encode %q", content)
	}
	resBits := new(utils.BitList)
	for i, r := range content {
		if i > 0 {
			resBits.AddBit(false)
		}
		resBits.AddBit(encodingTable[r]...)
	}
	return utils.New1DCodeWithColor(barcode.TypeCodabar, content, resBits, color), nil
}

// Encode creates a codabar barcode for the given content
func Encode(content string) (barcode.Barcode, error) {
	return EncodeWithColor(content, barcode.ColorScheme16)
}
