package qr

import (
	"bytes"
	"strings"
	"testing"
)

func Test_UnicodeEncoding(t *testing.T) {
	encode := Unicode.getEncoder()
	x, vi, err := encode("A", H) // 65
	if x == nil || vi == nil || vi.Version != 1 || !bytes.Equal(x.GetBytes(), []byte{64, 20, 16, 236, 17, 236, 17, 236, 17}) {
		t.Errorf("\"A\" failed to encode: %s", err)
	}
	_, _, err = encode(strings.Repeat("A", 3000), H)
	if err == nil {
		t.Error("Unicode encoding should not be able to encode a 3kb string")
	}
}
