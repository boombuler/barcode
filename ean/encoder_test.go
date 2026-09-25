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

func Test_EncodeChecksum(t *testing.T) {
	short8, err := Encode("5512345")
	if err != nil {
		t.Fatal(err)
	}
	full8, err := Encode("55123457")
	if err != nil {
		t.Fatal(err)
	}
	if short8.Content() != "55123457" || short8.CheckSum() != full8.CheckSum() {
		t.Errorf("EAN-8 short content %q checksum %d, full checksum %d", short8.Content(), short8.CheckSum(), full8.CheckSum())
	}

	short13, err := Encode("590123412345")
	if err != nil {
		t.Fatal(err)
	}
	full13, err := Encode("5901234123457")
	if err != nil {
		t.Fatal(err)
	}
	if short13.Content() != "5901234123457" || short13.CheckSum() != 7 || short13.CheckSum() != full13.CheckSum() {
		t.Errorf("EAN-13 short content %q checksum %d, full checksum %d", short13.Content(), short13.CheckSum(), full13.CheckSum())
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
