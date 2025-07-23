package qr

import (
	"bytes"
	"strings"
	"testing"
)

func Test_NumericEncoding(t *testing.T) {
	encode := Numeric.getEncoder()
	x, vi, err := encode("01234567", H)
	if err != nil {
		t.Fatal(err)
	}
	if x == nil || vi == nil || vi.Version != 1 || !bytes.Equal(x.GetBytes(), []byte{16, 32, 12, 86, 97, 128, 236, 17, 236}) {
		t.Error("\"01234567\" failed to encode")
	}
	x, vi, err = encode("0123456789012345", H)
	if err != nil {
		t.Fatal(err)
	}
	if x == nil || vi == nil || vi.Version != 1 || !bytes.Equal(x.GetBytes(), []byte{16, 64, 12, 86, 106, 110, 20, 234, 80}) {
		t.Error("\"0123456789012345\" failed to encode")
	}
	_, _, err = encode("foo", H)
	if err == nil {
		t.Error("Numeric encoding should not be able to encode \"foo\"")
	}
	x, vi, err = encode(strings.Repeat("1", 14297), H)
	if x != nil || vi != nil || err == nil {
		t.Fail()
	}
}
