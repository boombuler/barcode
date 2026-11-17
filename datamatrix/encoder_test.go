package datamatrix

import (
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/internal/imagetest"
)

func TestImage(t *testing.T) {
	imagetest.Run(t, imagetest.Encoder{
		Encode: func() (barcode.Barcode, error) {
			return Encode("1234567890")
		},
		EncodeWithColor: func(s barcode.ColorScheme) (barcode.Barcode, error) {
			return EncodeWithColor("1234567890", s)
		},
	})
}
