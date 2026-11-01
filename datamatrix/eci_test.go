package datamatrix

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/boombuler/barcode"
)

func Test_ECICodewords(t *testing.T) {
	tests := []struct {
		eci   barcode.ECI
		words []byte
	}{
		{barcode.NoECI, nil},
		{0, []byte{241, 1}},
		{barcode.ECIISO8859_1, []byte{241, 4}},
		{barcode.ECIUTF8, []byte{241, 27}},
		{126, []byte{241, 127}},
		{127, []byte{241, 128, 1}},
		{16382, []byte{241, 191, 254}},
		{16383, []byte{241, 192, 1, 1}},
		{barcode.MaxECI, []byte{241, 207, 63, 129}},
	}
	for _, test := range tests {
		if got := eciCodewords(test.eci); !bytes.Equal(got, test.words) {
			t.Errorf("eciCodewords(%d) = %v, want %v", test.eci, got, test.words)
		}
	}
}

func Test_ECIText(t *testing.T) {
	tests := []struct {
		content string
		eci     barcode.ECI
		words   []byte
	}{
		// "é" in ISO 8859-1: the ECI, Upper Shift, and the byte less 127.
		{"\xe9", barcode.ECIISO8859_1, []byte{241, 4, 235, 106}},
		// Without an ECI a leading FNC1 byte makes a GS1 symbol.
		{"\xe801", barcode.NoECI, []byte{232, 131}},
		// With one it is "è" in ISO 8859-1, as any other byte would be.
		{"\xe801", barcode.ECIISO8859_1, []byte{241, 4, 235, 105, 131}},
	}
	for _, test := range tests {
		got := encodeText(test.content, test.eci)
		if !bytes.Equal(got, test.words) {
			t.Errorf("encodeText(%q, %d) = %v, want %v",
				test.content, test.eci, got, test.words)
		}
	}
}

func Test_EncodeWithECI(t *testing.T) {
	content := "Caf\u00e9"
	plain, err := Encode(content)
	if err != nil {
		t.Fatal(err)
	}
	none, err := EncodeWithECI(content, barcode.NoECI)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plain.(*datamatrixCode).BitList,
		none.(*datamatrixCode).BitList) {
		t.Error("NoECI draws a different symbol from Encode")
	}
	withECI, err := EncodeWithECI(content, barcode.ECIUTF8)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(plain.(*datamatrixCode).BitList,
		withECI.(*datamatrixCode).BitList) {
		t.Error("an ECI draws the same symbol as no ECI")
	}
	if withECI.Content() != content {
		t.Errorf("Content() = %q, want %q", withECI.Content(), content)
	}

	for _, eci := range []barcode.ECI{-2, barcode.MaxECI + 1} {
		if _, err := EncodeWithECI(content, eci); err == nil {
			t.Errorf("ECI %d should be refused", eci)
		}
	}
}
