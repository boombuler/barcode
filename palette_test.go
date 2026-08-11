package barcode_test

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
	"github.com/boombuler/barcode/codabar"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/code93"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/pdf417"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
)

const (
	testWidth       = 900
	testHeight      = 300
	testPDF417Level = 3
)

type barcodeCase struct {
	name   string
	encode func(barcode.ColorScheme) (barcode.Barcode, error)
}

var barcodeCases = []barcodeCase{
	{
		name: "Aztec",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return aztec.EncodeWithColor([]byte("1234567890"), 23, aztec.DEFAULT_LAYERS, scheme)
		},
	},
	{
		name: "Codabar",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return codabar.EncodeWithColor("A12345B", scheme)
		},
	},
	{
		name: "Code128",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return code128.EncodeWithColor("1234567890123456", scheme)
		},
	},
	{
		name: "Code39",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return code39.EncodeWithColor("1234567890", false, false, scheme)
		},
	},
	{
		name: "Code93",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return code93.EncodeWithColor("1234567890", false, false, scheme)
		},
	},
	{
		name: "DataMatrix",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return datamatrix.EncodeWithColor("1234567890", scheme)
		},
	},
	{
		name: "EAN",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return ean.EncodeWithColor("5901234123457", scheme)
		},
	},
	{
		name: "PDF417",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return pdf417.EncodeWithColor("1234567890123456789012345678901234567890", testPDF417Level, scheme)
		},
	},
	{
		name: "QRCode",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return qr.EncodeWithColor("1234567890", qr.M, qr.Auto, scheme)
		},
	},
	{
		name: "TwoOfFive",
		encode: func(scheme barcode.ColorScheme) (barcode.Barcode, error) {
			return twooffive.EncodeWithColor("1234567890", false, scheme)
		},
	},
}

func BenchmarkGenerateBarcode(b *testing.B) {
	variants := []struct {
		name   string
		scheme barcode.ColorScheme
	}{
		{name: "gray16", scheme: barcode.ColorScheme16},
		{name: "palette", scheme: barcode.ColorSchemePalette},
	}

	for _, tc := range barcodeCases {
		for _, variant := range variants {
			b.Run(tc.name+"/"+variant.name, func(b *testing.B) {
				b.ReportAllocs()
				pngSize := 0

				for i := 0; i < b.N; i++ {
					code, err := tc.encode(variant.scheme)
					if err != nil {
						b.Fatal(err)
					}

					scaled, err := barcode.Scale(code, testWidth, testHeight)
					if err != nil {
						b.Fatal(err)
					}

					var output bytes.Buffer
					if err := png.Encode(&output, scaled); err != nil {
						b.Fatal(err)
					}
					pngSize = output.Len()
				}

				b.ReportMetric(float64(pngSize), "png_B")
			})
		}
	}
}
