// Package imagetest checks that barcodes implement [image.Image] and
// [image.PalettedImage] consistently with their color scheme, for the tests
// of each symbology.
package imagetest

import (
	"image"
	"image/color"
	"testing"

	"github.com/boombuler/barcode"
)

// Scheme is a color scheme to encode a barcode with.
type Scheme struct {
	Name string
	barcode.ColorScheme
}

// Default holds the colors Encode is expected to use. Its Model is unset, as
// the default color model is not checked.
var Default = barcode.ColorScheme{
	Background: color.White,
	Foreground: color.Black,
}

// Schemes cover each way that [image/png.Encode] reads the image.
var Schemes = []Scheme{
	{"palette", barcode.ColorScheme{
		Model:      color.Palette{color.White, color.Black},
		Background: color.White,
		Foreground: color.Black,
	}},
	// Foreground and background should not depend on the order in the palette.
	{"reversed palette", barcode.ColorScheme{
		Model:      color.Palette{color.Black, color.White},
		Background: color.White,
		Foreground: color.Black,
	}},
	{"blue on yellow", barcode.ColorScheme{
		Model:      color.RGBAModel,
		Background: color.RGBA{R: 255, G: 255, A: 255},
		Foreground: color.RGBA{B: 255, A: 255},
	}},
}

var red = color.RGBA{R: 255, A: 255}

// Encoder encodes the same content with the Encode and EncodeWithColor
// functions of a symbology.
type Encoder struct {
	Encode          func() (barcode.Barcode, error)
	EncodeWithColor func(barcode.ColorScheme) (barcode.Barcode, error)
}

// Run checks that barcodes produced by [Encoder] correctly implement
// the [image.Image] and [image.PalettedImage] for all color schemes
// and when [barcode.Scale] or [barcode.ScaleWithFill] is applied.
func Run(t *testing.T, e Encoder) {
	t.Helper()
	if e.Encode == nil || e.EncodeWithColor == nil {
		t.Fatal("imagetest: Encoder needs Encode and EncodeWithColor")
	}
	t.Run("Encode", func(t *testing.T) {
		code, err := e.Encode()
		if err != nil {
			t.Fatal(err)
		}
		check(t, code, Default)
	})
	for _, s := range Schemes {
		t.Run(s.Name, func(t *testing.T) {
			code, err := e.EncodeWithColor(s.ColorScheme)
			if err != nil {
				t.Fatal(err)
			}
			check(t, code, s.ColorScheme)
		})
	}
}

// check checks code raw and scaled.
func check(t *testing.T, code barcode.Barcode, scheme barcode.ColorScheme) {
	t.Run("raw", func(t *testing.T) {
		checkImage(t, code, code, scheme, nil)
	})

	// Ensure there is one pixel more on each side so fill can be tested.
	b := code.Bounds()
	width, height := b.Dx()*2+2, b.Dy()*2+2
	t.Run("Scale", func(t *testing.T) {
		img, err := barcode.Scale(code, width, height)
		if err != nil {
			t.Fatal(err)
		}
		checkImage(t, code, img, scheme, scheme.Background)
	})
	for _, fill := range []struct {
		name  string
		color color.Color
	}{
		{"red", red},
		{"transparent", color.Transparent},
	} {
		t.Run("ScaleWithFill "+fill.name, func(t *testing.T) {
			img, err := barcode.ScaleWithFill(code, width, height, fill.color)
			if err != nil {
				t.Fatal(err)
			}
			checkImage(t, code, img, scheme, fill.color)
		})
	}
}

// checkImage checks img, which is code or code scaled. fill is nil for code.
func checkImage(t *testing.T, code, img barcode.Barcode, scheme barcode.ColorScheme, fill color.Color) {
	t.Helper()
	b := img.Bounds()
	// The fill pads a scaled barcode.
	if fill != nil {
		if got := img.At(b.Min.X, b.Min.Y); got != fill {
			t.Errorf("got %#v in the corner, want the fill %#v", got, fill)
		}
	}

	model := code.ColorModel()
	palette, paletted := model.(color.Palette)
	var indexed image.PalettedImage
	if paletted {
		var ok bool
		if indexed, ok = img.(image.PalettedImage); !ok {
			t.Fatalf("%T does not implement image.PalettedImage", img)
		}
	}

	var hasBackground, hasForeground bool
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// Callers of the package are known to compare At with ==
			// so ensure given colors from scheme are returned.
			c := img.At(x, y)
			switch {
			case c == scheme.Background:
				hasBackground = true
			case c == scheme.Foreground:
				hasForeground = true
			case c != fill:
				t.Fatalf("got %#v at (%d, %d), want %#v, %#v or the fill %#v",
					c, x, y, scheme.Background, scheme.Foreground, fill)
			}

			// Encoders convert through ColorModel, so a scaled barcode must match conversion of original.
			if got, want := img.ColorModel().Convert(c), model.Convert(c); !sameColor(got, want) {
				t.Fatalf("ColorModel converts %#v at (%d, %d) to %#v, want %#v", c, x, y, got, want)
			}

			// Encoders of paletted images read ColorIndexAt instead.
			if paletted {
				i := int(indexed.ColorIndexAt(x, y))
				if i >= len(palette) {
					t.Fatalf("got index %d at (%d, %d), want less than %d", i, x, y, len(palette))
				}
				if got, want := palette[i], palette.Convert(c); got != want {
					t.Fatalf("got index %d (%#v) at (%d, %d), want %d (%#v)",
						i, got, x, y, palette.Index(c), want)
				}
			}
		}
	}
	if !hasBackground || !hasForeground {
		t.Errorf("got background %t and foreground %t, want both", hasBackground, hasForeground)
	}
}

func sameColor(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
