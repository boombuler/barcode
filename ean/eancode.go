package ean

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
	"image"
	"image/color"
)

type eancode struct {
	*utils.BitList
	content string
}

func newEANCode(isEAN8 bool) *eancode {
	capacity := 95
	if isEAN8 {
		capacity = 67
	}
	return &eancode{utils.NewBitList(capacity), ""}
}

func (c *eancode) Content() string {
	return c.content
}

func (c *eancode) Metadata() barcode.Metadata {
	if c.Len() == 67 {
		return barcode.Metadata{"EAN 8", 1}
	}
	return barcode.Metadata{"EAN 13", 1}
}

func (c *eancode) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *eancode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Len(), 1)
}

func (c *eancode) At(x, y int) color.Color {
	if c.GetBit(x) {
		return color.Black
	}
	return color.White
}
