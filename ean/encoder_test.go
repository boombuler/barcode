package ean

import (
	"image/color"
	"testing"
)

func Test_EncodeEAN13(t *testing.T) {
	testResult := "10100010110100111011001100100110111101001110101010110011011011001000010101110010011101000100101"
	testCode := "5901234123457"
	code, err := Encode(testCode)
	if err != nil {
		t.Error(err)
	}
	if code.Metadata().Dimensions != 1 || code.Content() != testCode || code.Metadata().CodeKind != "EAN 13" {
		t.Error("Metadata missmatch")
	}
	for i, r := range testResult {
		if (code.At(i, 0) == color.Black) != (r == '1') {
			t.Fail()
		}
	}
}
