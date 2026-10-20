[![Join the chat at https://gitter.im/golang-barcode/Lobby](https://badges.gitter.im/golang-barcode/Lobby.svg)](https://gitter.im/golang-barcode/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Introduction ##

This is a package for GO which can be used to create different types of barcodes.

## Supported Barcode Types ##
* 2 of 5
* Aztec Code
* Codabar
* Code 128
* Code 39
* Code 93
* Datamatrix
* EAN 13
* EAN 8
* PDF 417
* QR Code

## Example ##

This is a simple example on how to create a QR-Code and write it to a png-file
```go
package main

import (
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	// Create the barcode
	qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	// create the output file
	file, _ := os.Create("qrcode.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)
}
```

## Documentation ##
See [GoDoc](https://godoc.org/github.com/boombuler/barcode)

To create a barcode use the Encode function from one of the subpackages.

### Scaling with a fill color ###

`ScaleWithFill` adds padding using the requested fill but keeps the barcode's
original color model. A fill outside that model can change color or lose
transparency when the image is encoded. For example, the default QR code uses
grayscale (`ColorScheme16`), so transparent padding becomes opaque in a PNG.

Use `qr.EncodeWithColor` with `barcode.ColorScheme32` before calling
`ScaleWithFill` when you need colored or transparent padding. This changes the
color model; the barcode's background remains white unless you also customize
the color scheme. See the executable
[`ScaleWithFill` example](https://pkg.go.dev/github.com/boombuler/barcode#example-ScaleWithFill).
