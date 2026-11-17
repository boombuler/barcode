package barcode_test

import (
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// A transparent fill needs a color model that holds it, so the QR code is
// encoded with ColorScheme32 rather than the default grayscale scheme.
func ExampleScaleWithFill() {
	code, err := qr.EncodeWithColor("Hello World", qr.M, qr.Auto, barcode.ColorScheme32)
	if err != nil {
		log.Fatal(err)
	}

	code, err = barcode.ScaleWithFill(code, 300, 200, color.Transparent)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, code); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
