package barcode

import "image/color"

// ColorScheme defines a structure for color schemes used in barcode rendering.
// It includes the color model, background color, and foreground color.
type ColorScheme struct {
	// Color model to be used (e.g., grayscale, RGB, RGBA)
	//
	// A color.Palette makes the barcodes render as paletted images. Such a
	// palette must list Background as its first and Foreground as its second
	// entry, as that is the order the barcodes report their palette indexes in.
	// Use NewPaletteColorScheme to build a palette which holds to that order.
	Model      color.Model
	Background color.Color // Color of the background
	Foreground color.Color // Color of the foreground (e.g., bars in a barcode)
}

// ColorScheme8 represents a color scheme with 8-bit grayscale colors.
var ColorScheme8 = ColorScheme{
	Model:      color.GrayModel,
	Background: color.Gray{Y: 255},
	Foreground: color.Gray{Y: 0},
}

// ColorScheme16 represents a color scheme with 16-bit grayscale colors.
var ColorScheme16 = ColorScheme{
	Model:      color.Gray16Model,
	Background: color.White,
	Foreground: color.Black,
}

// ColorScheme24 represents a color scheme with 24-bit RGB colors.
var ColorScheme24 = ColorScheme{
	Model:      color.RGBAModel,
	Background: color.RGBA{255, 255, 255, 255},
	Foreground: color.RGBA{0, 0, 0, 255},
}

// ColorScheme32 represents a color scheme with 32-bit RGBA colors, which is similar to ColorScheme24 but typically includes alpha for transparency.
var ColorScheme32 = ColorScheme{
	Model:      color.RGBAModel,
	Background: color.RGBA{255, 255, 255, 255},
	Foreground: color.RGBA{0, 0, 0, 255},
}

// NewPaletteColorScheme returns a color scheme which renders paletted images
// using the given background and foreground colors, with the palette in the
// order the barcodes report their color indexes in.
func NewPaletteColorScheme(background, foreground color.Color) ColorScheme {
	return ColorScheme{
		Model:      color.Palette{background, foreground},
		Background: background,
		Foreground: foreground,
	}
}

// ColorSchemePalette represents a black on white color scheme that uses a
// palette of colors. Encoders which support paletted images store a palette
// index per pixel instead of a full color, so a two color barcode needs a
// single bit per pixel.
var ColorSchemePalette = NewPaletteColorScheme(color.White, color.Black)
