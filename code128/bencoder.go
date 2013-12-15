package code128

import (
	"github.com/boombuler/barcode/utils"
	"strings"
)

const bTable = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

func encodeBTable(content []rune) *utils.BitList {
	result := new(utils.BitList)
	result.AddByte(startBSymbol)
	for _, r := range content {
		idx := strings.IndexRune(bTable, r)
		if idx < 0 {
			return nil
		}
		result.AddByte(byte(idx))
	}
	return result
}
