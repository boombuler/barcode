package qr

import (
	"bytes"
	"strings"
	"testing"
)

func Test_AlphaNumericEncoding(t *testing.T) {
	encode := AlphaNumeric.getEncoder()

	x, vi, err := encode("HELLO WORLD", M)

	if x == nil || vi == nil || vi.Version != 1 || !bytes.Equal(x.GetBytes(), []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236, 17, 236, 17}) {
		t.Errorf("\"HELLO WORLD\" failed to encode: %s", err)
	}

	x, vi, err = encode(strings.Repeat("A", 4296), L)
	if x == nil || vi == nil || err != nil {
		t.Fail()
	}
	x, vi, err = encode(strings.Repeat("A", 4297), L)
	if x != nil || vi != nil || err == nil {
		t.Fail()
	}
	x, vi, err = encode("ABc", L)
	if x != nil || vi != nil || err == nil {
		t.Fail()
	}
	x, vi, err = encode("hello world", M)

	if x != nil || vi != nil || err == nil {
		t.Error("\"hello world\" should not be encodable in alphanumeric mode")
	}
}
