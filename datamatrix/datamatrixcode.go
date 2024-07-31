package datamatrix

import (
	"image"
	"image/color"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type datamatrixCode struct {
	*utils.BitList
	*dmCodeSize
	content string
	color   barcode.ColorScheme
}

func newDataMatrixCodeWithColor(size *dmCodeSize, color barcode.ColorScheme) *datamatrixCode {
	return &datamatrixCode{utils.NewBitList(size.Rows * size.Columns), size, "", color}
}

func newDataMatrixCode(size *dmCodeSize) *datamatrixCode {
	return &datamatrixCode{utils.NewBitList(size.Rows * size.Columns), size, "", barcode.ColorScheme16}
}

func (c *datamatrixCode) Content() string {
	return c.content
}

func (c *datamatrixCode) Metadata() barcode.Metadata {
	return barcode.Metadata{barcode.TypeDataMatrix, 2}
}

func (c *datamatrixCode) ColorModel() color.Model {
	return c.color.Model
}

func (c *datamatrixCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Columns, c.Rows)
}

func (c *datamatrixCode) At(x, y int) color.Color {
	if c.get(x, y) {
		return c.color.Foreground
	}
	return c.color.Background
}

func (c *datamatrixCode) get(x, y int) bool {
	return c.GetBit(x*c.Rows + y)
}

func (c *datamatrixCode) set(x, y int, value bool) {
	c.SetBit(x*c.Rows+y, value)
}
