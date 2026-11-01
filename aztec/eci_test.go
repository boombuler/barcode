package aztec

import (
	"reflect"
	"strings"
	"testing"

	"github.com/boombuler/barcode"
)

func Test_HighLevelEncodeECI(t *testing.T) {
	tests := []struct {
		data string
		eci  barcode.ECI
		bits string
	}{
		// P/S   FLG(n) n    '2'  '6'
		{"", barcode.ECIUTF8,
			"..... ..... .X. .X.. X..."},
		// P/S   FLG(n) n    '2'  '6'  'A'
		{"A", barcode.ECIUTF8,
			"..... ..... .X. .X.. X... ...X."},
		// P/S   FLG(n) n    '3'  L/L   'a'
		{"a", barcode.ECIISO8859_1,
			"..... ..... ..X .X.X XXX.. ...X."},
		// P/S   FLG(n) n    '0'
		{"", 0,
			"..... ..... ..X ..X."},
		// P/S   FLG(n) n    '9'  '9'  '9'  '9'  '9'  '9'
		{"", barcode.MaxECI,
			"..... ..... XX. X.XX X.XX X.XX X.XX X.XX X.XX"},
	}
	for _, test := range tests {
		want := strings.ReplaceAll(test.bits, " ", "")
		got := bitStr(highlevelEncode([]byte(test.data), test.eci))
		if got != want {
			t.Errorf("highlevelEncode(%q, %d) = %s, want %s",
				test.data, test.eci, got, want)
		}
	}
}

func Test_EncodeWithECI(t *testing.T) {
	data := []byte("Caf\u00e9")
	encode := func(eci barcode.ECI) (barcode.Barcode, error) {
		return EncodeWithECI(data, DEFAULT_EC_PERCENT, DEFAULT_LAYERS, eci)
	}
	plain, err := Encode(data, DEFAULT_EC_PERCENT, DEFAULT_LAYERS)
	if err != nil {
		t.Fatal(err)
	}
	none, err := encode(barcode.NoECI)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plain.(*aztecCode).BitList,
		none.(*aztecCode).BitList) {
		t.Error("NoECI draws a different symbol from Encode")
	}
	withECI, err := encode(barcode.ECIUTF8)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(plain.(*aztecCode).BitList,
		withECI.(*aztecCode).BitList) {
		t.Error("an ECI draws the same symbol as no ECI")
	}

	for _, eci := range []barcode.ECI{-2, barcode.MaxECI + 1} {
		if _, err := encode(eci); err == nil {
			t.Errorf("ECI %d should be refused", eci)
		}
	}
}
