package code39

import (
	"github.com/boombuler/barcode"
	"image"
	"image/color"
)

type code struct {
	*barcode.BitList
	content string
}

func newCode() *code {
	return &code{&barcode.BitList{}, ""}
}

func (c *code) Content() string {
	return c.content
}

func (c *code) Metadata() barcode.Metadata {
	return barcode.Metadata{"Code39", 1}
}

func (c *code) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *code) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Len(), 1)
}

func (c *code) At(x, y int) color.Color {
	if c.GetBit(x) {
		return color.Black
	}
	return color.White
}
