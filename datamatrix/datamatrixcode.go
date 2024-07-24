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
	depth   int
}

func newDataMatrixCodeWithDepth(size *dmCodeSize, depth int) *datamatrixCode {
	return &datamatrixCode{utils.NewBitList(size.Rows * size.Columns), size, "", depth}
}

func newDataMatrixCode(size *dmCodeSize) *datamatrixCode {
	return &datamatrixCode{utils.NewBitList(size.Rows * size.Columns), size, "", 16}
}

func (c *datamatrixCode) Content() string {
	return c.content
}

func (c *datamatrixCode) Metadata() barcode.Metadata {
	return barcode.Metadata{barcode.TypeDataMatrix, 2}
}

func (c *datamatrixCode) ColorModel() color.Model {
	return utils.ColorModel(c.depth)
}

func (c *datamatrixCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Columns, c.Rows)
}

func (c *datamatrixCode) At(x, y int) color.Color {
	if c.get(x, y) {
		return utils.BlackColor(c.depth)
	}
	return utils.WhiteColor(c.depth)
}

func (c *datamatrixCode) get(x, y int) bool {
	return c.GetBit(x*c.Rows + y)
}

func (c *datamatrixCode) set(x, y int, value bool) {
	c.SetBit(x*c.Rows+y, value)
}
