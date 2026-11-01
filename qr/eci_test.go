package qr

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

func bitString(bits *utils.BitList) string {
	var buf strings.Builder
	for idx := 0; idx < bits.Len(); idx++ {
		if bits.GetBit(idx) {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	}
	return buf.String()
}

func Test_ECIHeader(t *testing.T) {
	tests := []struct {
		eci  barcode.ECI
		bits string
	}{
		{barcode.NoECI, ""},
		{0, "0111 00000000"},
		{barcode.ECIISO8859_1, "0111 00000011"},
		{barcode.ECIUTF8, "0111 00011010"},
		{127, "0111 01111111"},
		{128, "0111 10 00000010000000"},
		{16383, "0111 10 11111111111111"},
		{16384, "0111 110 000000100000000000000"},
		{barcode.MaxECI, "0111 110 011110100001000111111"},
	}
	for _, test := range tests {
		want := strings.ReplaceAll(test.bits, " ", "")
		if got := bitString(eciHeader(test.eci)); got != want {
			t.Errorf("eciHeader(%d) = %s, want %s", test.eci, got, want)
		}
	}
}

func Test_ECIEncoding(t *testing.T) {
	// The ECI header, then "é" in UTF-8 in byte mode,
	// the terminator, and the padding.
	bits, vi, err := Unicode.getEncoder()("\u00e9", H, barcode.ECIUTF8)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x71, 0xA4, 0x02, 0xC3, 0xA9, 0x00, 0xEC, 0x11, 0xEC}
	if vi.Version != 1 || !bytes.Equal(bits.GetBytes(), want) {
		t.Errorf("version %d, bytes %x, want version 1, bytes %x",
			vi.Version, bits.GetBytes(), want)
	}

	// Auto passes the ECI to the mode it chooses.
	bits, _, err = Auto.getEncoder()("123", M, barcode.ECIUTF8)
	if err != nil {
		t.Fatal(err)
	}
	if got := bitString(bits); !strings.HasPrefix(got, "0111000110100001") {
		t.Errorf("Auto with an ECI starts %s, want the ECI, then numeric mode",
			got[:16])
	}

	// In UTF-16BE the bytes "33" are one character, U+3333, not two digits.
	bits, _, err = Auto.getEncoder()("33", M, barcode.ECIUTF16BE)
	if err != nil {
		t.Fatal(err)
	}
	if got := bitString(bits); !strings.HasPrefix(got, "0111000110010100") {
		t.Errorf("Auto under UTF-16BE starts %s, want the ECI, then byte mode",
			got[:16])
	}
}

func Test_ASCIICompatible(t *testing.T) {
	tests := map[barcode.ECI]bool{
		barcode.NoECI:        true,
		0:                    true,
		barcode.ECIISO8859_1: true,
		14:                   false,
		barcode.ECIShiftJIS:  true,
		barcode.ECIUTF16BE:   false,
		barcode.ECIUTF8:      true,
		barcode.ECIGB18030:   true,
		barcode.ECIUTF16LE:   false,
		barcode.ECIUTF32LE:   false,
		899:                  false,
	}
	for eci, want := range tests {
		if got := asciiCompatible(eci); got != want {
			t.Errorf("asciiCompatible(%d) = %v, want %v", eci, got, want)
		}
	}
}

func Test_ECIVersion(t *testing.T) {
	// 17 bytes fill version 1 at L so that the header no longer fits.
	content := strings.Repeat("a", 17)
	_, vi, err := Unicode.getEncoder()(content, L, barcode.NoECI)
	if err != nil || vi.Version != 1 {
		t.Fatalf("without an ECI: version %v, error %v, want version 1", vi, err)
	}
	_, vi, err = Unicode.getEncoder()(content, L, barcode.ECIUTF8)
	if err != nil || vi.Version != 2 {
		t.Fatalf("with an ECI: version %v, error %v, want version 2", vi, err)
	}
}

func Test_EncodeWithECI(t *testing.T) {
	content := "Caf\u00e9"
	plain, err := Encode(content, M, Auto)
	if err != nil {
		t.Fatal(err)
	}
	none, err := EncodeWithECI(content, M, Auto, barcode.NoECI)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plain.(*qrcode).data, none.(*qrcode).data) {
		t.Error("NoECI draws a different symbol from Encode")
	}
	withECI, err := EncodeWithECI(content, M, Auto, barcode.ECIUTF8)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(plain.(*qrcode).data, withECI.(*qrcode).data) {
		t.Error("an ECI draws the same symbol as no ECI")
	}
	if withECI.Content() != content {
		t.Errorf("Content() = %q, want %q", withECI.Content(), content)
	}

	for _, eci := range []barcode.ECI{-2, barcode.MaxECI + 1} {
		if _, err := EncodeWithECI(content, M, Auto, eci); err == nil {
			t.Errorf("ECI %d should be refused", eci)
		}
	}
	_, err = EncodeWithECI(content, M, unknownEncoding, barcode.NoECI)
	if err == nil {
		t.Error("an unknown encoding should be refused")
	}
}
