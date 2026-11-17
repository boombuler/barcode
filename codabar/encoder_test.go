package codabar

import (
	"image/color"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/internal/imagetest"
)

func Test_Encode(t *testing.T) {
	_, err := Encode("FOOBAR")
	if err == nil {
		t.Error("\"FOOBAR\" should not be encodable")
	}

	testEncode := func(txt, testResult string) {
		code, err := Encode(txt)
		if err != nil || code == nil {
			t.Fail()
		} else {
			if code.Bounds().Max.X != len(testResult) {
				t.Errorf("%v: length missmatch", txt)
			} else {
				for i, r := range testResult {
					if (code.At(i, 0) == color.Black) != (r == '1') {
						t.Errorf("%v: code missmatch on position %d", txt, i)
					}
				}
			}
		}
	}

	testEncode("A40156B", "10110010010101101001010101001101010110010110101001010010101101001001011")
}

func TestImage(t *testing.T) {
	imagetest.Run(t, imagetest.Encoder{
		Encode: func() (barcode.Barcode, error) {
			return Encode("A12345B")
		},
		EncodeWithColor: func(s barcode.ColorScheme) (barcode.Barcode, error) {
			return EncodeWithColor("A12345B", s)
		},
	})
}
