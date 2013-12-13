package codabar

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
	"image"
	"image/color"
)

type codabarcode struct {
	*utils.BitList
	content string
}

func (c *codabarcode) Content() string {
	return c.content
}

func (c *codabarcode) Metadata() barcode.Metadata {
	return barcode.Metadata{"Codabar", 1}
}

func (c *codabarcode) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *codabarcode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Len(), 1)
}

func (c *codabarcode) At(x, y int) color.Color {
	if c.GetBit(x) {
		return color.Black
	}
	return color.White
}
