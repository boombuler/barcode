package pdf417

import (
	"image"
	"image/color"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

type pdfBarcode struct {
	data  string
	width int
	code  *utils.BitList
	color barcode.ColorScheme
}

func (c *pdfBarcode) Metadata() barcode.Metadata {
	return barcode.Metadata{CodeKind: barcode.TypePDF, Dimensions: 2}
}

func (c *pdfBarcode) Content() string {
	return c.data
}

func (c *pdfBarcode) ColorModel() color.Model {
	return c.color.Model
}

func (c *pdfBarcode) ColorScheme() barcode.ColorScheme {
	return c.color
}

func (c *pdfBarcode) Bounds() image.Rectangle {
	height := c.code.Len() / c.width

	return image.Rect(0, 0, c.width, height*moduleHeight)
}

func (c *pdfBarcode) At(x, y int) color.Color {
	if c.code.GetBit((y/moduleHeight)*c.width + x) {
		return c.color.Foreground
	}
	return c.color.Background
}
