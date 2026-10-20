package barcode_test

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func ExampleScaleWithFill() {
	// Select an alpha-capable color model before adding transparent padding.
	code, err := qr.EncodeWithColor("Hello World", qr.M, qr.Auto, barcode.ColorScheme32)
	if err != nil {
		panic(err)
	}
	scaled, err := barcode.ScaleWithFill(code, 300, 200, color.Transparent)
	if err != nil {
		panic(err)
	}

	// Check the encoded image: At alone does not reflect color-model conversion.
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, scaled); err != nil {
		panic(err)
	}
	image, err := png.Decode(&buffer)
	if err != nil {
		panic(err)
	}
	_, _, _, alpha := image.At(0, 0).RGBA()
	fmt.Println("Padding alpha:", alpha)
	// Output: Padding alpha: 0
}
