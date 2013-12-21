package codabar

import (
	"image/color"
	"testing"
)

func Test_Encode(t *testing.T) {
	_, err := Encode("FOOBAR")
	if err == nil {
		t.Error("\"FOOBAR\" should not be encodable")
	}

	code, err := Encode("A40156B")
	if err != nil || code == nil {
		t.Fail()
	} else {
		testResult := "10110010010101101001010101001101010110010110101001010010101101010010011"
		if code.Bounds().Max.X != len(testResult) {
			t.Error("length missmatch")
		} else {
			for i, r := range testResult {
				if (code.At(i, 0) == color.Black) != (r == '1') {
					t.Errorf("code missmatch on position %d", i)
				}
			}
		}
	}
}
