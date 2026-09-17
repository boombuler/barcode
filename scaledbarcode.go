package barcode

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
)

type wrapFunc func(x, y int) color.Color

type scaledBarcode struct {
	wrapped     Barcode
	wrapperFunc wrapFunc
	rect        image.Rectangle
	// palette of the wrapped barcode, nil if it does not render paletted images
	palette color.Palette
}

type intCSscaledBC struct {
	scaledBarcode
}

var (
	_ image.PalettedImage = (*scaledBarcode)(nil)
	_ image.PalettedImage = (*intCSscaledBC)(nil)
)

func (bc *scaledBarcode) Content() string {
	return bc.wrapped.Content()
}

func (bc *scaledBarcode) Metadata() Metadata {
	return bc.wrapped.Metadata()
}

func (bc *scaledBarcode) ColorModel() color.Model {
	return bc.wrapped.ColorModel()
}

func (bc *scaledBarcode) Bounds() image.Rectangle {
	return bc.rect
}

func (bc *scaledBarcode) At(x, y int) color.Color {
	return bc.wrapperFunc(x, y)
}

// ColorIndexAt looks the color of a pixel up in the palette of the wrapped
// barcode. Resolving it through At keeps the index in sync with the rendered
// color and covers the fill around the barcode, which is not required to be
// part of the palette.
func (bc *scaledBarcode) ColorIndexAt(x, y int) uint8 {
	// palettes larger than 256 entries can not be addressed by a uint8 index
	// and are not renderable, so those fall back to the first entry.
	if idx := bc.palette.Index(bc.At(x, y)); idx >= 0 && idx <= 255 {
		return uint8(idx)
	}
	return 0
}

func (bc *intCSscaledBC) CheckSum() int {
	if cs, ok := bc.wrapped.(BarcodeIntCS); ok {
		return cs.CheckSum()
	}
	return 0
}

// Scale returns a resized barcode with the given width and height.
func Scale(bc Barcode, width, height int) (Barcode, error) {
	var fill color.Color
	if v, ok := bc.(BarcodeColor); ok {
		fill = v.ColorScheme().Background
	} else {
		fill = color.White
	}
	return ScaleWithFill(bc, width, height, fill)
}

// Scale returns a resized barcode with the given width, height and fill color.
func ScaleWithFill(bc Barcode, width, height int, fill color.Color) (Barcode, error) {
	switch bc.Metadata().Dimensions {
	case 1:
		return scale1DCode(bc, width, height, fill)
	case 2:
		return scale2DCode(bc, width, height, fill)
	}

	return nil, errors.New("unsupported barcode format")
}

func newScaledBC(wrapped Barcode, wrapperFunc wrapFunc, rect image.Rectangle) Barcode {
	result := &scaledBarcode{
		wrapped:     wrapped,
		wrapperFunc: wrapperFunc,
		rect:        rect,
	}
	result.palette, _ = wrapped.ColorModel().(color.Palette)

	if _, ok := wrapped.(BarcodeIntCS); ok {
		return &intCSscaledBC{*result}
	}
	return result
}

func scale2DCode(bc Barcode, width, height int, fill color.Color) (Barcode, error) {
	orgBounds := bc.Bounds()
	orgWidth := orgBounds.Max.X - orgBounds.Min.X
	orgHeight := orgBounds.Max.Y - orgBounds.Min.Y

	factor := int(math.Min(float64(width)/float64(orgWidth), float64(height)/float64(orgHeight)))
	if factor <= 0 {
		return nil, fmt.Errorf("can not scale barcode to an image smaller than %dx%d", orgWidth, orgHeight)
	}

	offsetX := (width - (orgWidth * factor)) / 2
	offsetY := (height - (orgHeight * factor)) / 2

	wrap := func(x, y int) color.Color {
		if x < offsetX || y < offsetY {
			return fill
		}
		x = (x - offsetX) / factor
		y = (y - offsetY) / factor
		if x >= orgWidth || y >= orgHeight {
			return fill
		}
		return bc.At(x, y)
	}

	return newScaledBC(
		bc,
		wrap,
		image.Rect(0, 0, width, height),
	), nil
}

func scale1DCode(bc Barcode, width, height int, fill color.Color) (Barcode, error) {
	orgBounds := bc.Bounds()
	orgWidth := orgBounds.Max.X - orgBounds.Min.X
	factor := int(float64(width) / float64(orgWidth))

	if factor <= 0 {
		return nil, fmt.Errorf("can not scale barcode to an image smaller than %dx1", orgWidth)
	}
	offsetX := (width - (orgWidth * factor)) / 2

	wrap := func(x, y int) color.Color {
		if x < offsetX {
			return fill
		}
		x = (x - offsetX) / factor

		if x >= orgWidth {
			return fill
		}
		return bc.At(x, 0)
	}

	return newScaledBC(
		bc,
		wrap,
		image.Rect(0, 0, width, height),
	), nil
}
