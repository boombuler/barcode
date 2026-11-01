package qr

import (
	"bytes"
	"strings"
	"testing"

	"github.com/boombuler/barcode"
)

func Test_AutomaticEncoding(t *testing.T) {
	tests := map[string]encodeFn{
		"0123456789":      Numeric.getEncoder(),
		"ALPHA NUMERIC":   AlphaNumeric.getEncoder(),
		"unicode encoing": Unicode.getEncoder(),
		"very long unicode encoding" + strings.Repeat("A", 3000): nil,
	}

	for str, enc := range tests {
		testValue, _, _ := Auto.getEncoder()(str, M, barcode.NoECI)
		if enc != nil {
			correctValue, _, _ := enc(str, M, barcode.NoECI)
			if testValue == nil || !bytes.Equal(correctValue.GetBytes(), testValue.GetBytes()) {
				t.Errorf("wrong encoding used for '%s'", str)
			}
		} else {
			if testValue != nil {
				t.Errorf("wrong encoding used for '%s'", str)
			}
		}

	}
}
