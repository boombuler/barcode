package qr

import (
	"bytes"
	"testing"
)

func makeString(length int) string {
	res := ""

	for i := 0; i < length; i++ {
		res += "A"
	}

	return res
}

func Test_AlphaNumericEncoding(t *testing.T) {
	x, vi, err := AlphaNumeric.encode("HELLO WORLD", M)

	if x == nil || vi == nil || vi.Version != 1 || bytes.Compare(x.GetBytes(), []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236, 17, 236, 17}) != 0 {
		t.Errorf("\"HELLO WORLD\" failed to encode: %s", err)
	}

	x, vi, err = AlphaNumeric.encode(makeString(4296), L)
	if x == nil || vi == nil || err != nil {
		t.Fail()
	}
	x, vi, err = AlphaNumeric.encode(makeString(4297), L)
	if x != nil || vi != nil || err == nil {
		t.Fail()
	}
	x, vi, err = AlphaNumeric.encode("ABc", L)
	if x != nil || vi != nil || err == nil {
		t.Fail()
	}
	x, vi, err = AlphaNumeric.encode("hello world", M)

	if x != nil || vi != nil || err == nil {
		t.Error("\"hello world\" should not be encodable in alphanumeric mode")
	}
}
