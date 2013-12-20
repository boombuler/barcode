package ean

import (
	"image/color"
	"testing"
)

func testHelper(t *testing.T, testCode, testResult, kind string, checkMetadata bool) {
	code, err := Encode(testCode)
	if err != nil {
		t.Error(err)
	}
	if checkMetadata && (code.Metadata().Dimensions != 1 || code.Content() != testCode || code.Metadata().CodeKind != kind) {
		t.Error("Metadata missmatch")
	}
	if len(testResult) != code.Bounds().Max.X {
		t.Fail()
	}
	for i, r := range testResult {
		if (code.At(i, 0) == color.Black) != (r == '1') {
			t.Fail()
		}
	}
}

func Test_EncodeEAN(t *testing.T) {
	testHelper(t, "5901234123457", "10100010110100111011001100100110111101001110101010110011011011001000010101110010011101000100101", "EAN 13", true)
	testHelper(t, "55123457", "1010110001011000100110010010011010101000010101110010011101000100101", "EAN 8", true)
	testHelper(t, "5512345", "1010110001011000100110010010011010101000010101110010011101000100101", "EAN 8", false)
	_, err := Encode("55123458") //<-- Invalid checksum
	if err == nil {
		t.Error("Invalid checksum not detected")
	}
	_, err = Encode("invalid")
	if err == nil {
		t.Error("\"invalid\" should not be encodable")
	}
	_, err = Encode("invalid")
	if err == nil {
		t.Error("\"invalid\" should not be encodable")
	}
	bits := encodeEAN13("invalid error")
	if bits != nil {
		t.Error("\"invalid error\" should not be encodable")
	}
}

func Test_RuneToIntIntToRune(t *testing.T) {
	if intToRune(0) != '0' {
		t.Errorf("failed intToRune(0) returned %d", string(intToRune(0)))
	}
	if intToRune(9) != '9' {
		t.Errorf("failed intToRune(9) returned %d", intToRune(9))
	}
	if intToRune(10) != 'F' {
		t.Errorf("failed intToRune(10) returned %d", intToRune(10))
	}
	if runeToInt('0') != 0 {
		t.Error("failed runeToInt('0') returned %d", runeToInt(0))
	}
	if runeToInt('9') != 9 {
		t.Error("failed runeToInt('9') returned %d", runeToInt(9))
	}
	if runeToInt('F') != -1 {
		t.Error("failed runeToInt('F') returned %d", runeToInt('F'))
	}
}
