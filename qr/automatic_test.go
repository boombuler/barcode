package qr

import (
	"bytes"
	"testing"
)

func Test_AutomaticEncoding(t *testing.T) {
	tests := map[string]Encoding{
		"0123456789":          Numeric,
		"ALPHA NUMERIC":       AlphaNumeric,
		"no matching encoing": nil,
	}

	for str, enc := range tests {
		testValue, _, _ := Auto.encode(str, M)
		if enc != nil {
			correctValue, _, _ := enc.encode(str, M)
			if testValue == nil || bytes.Compare(correctValue.GetBytes(), testValue.GetBytes()) != 0 {
				t.Errorf("wrong encoding used for '%s'", str)
			}
		} else {
			if testValue != nil {
				t.Errorf("wrong encoding used for '%s'", str)
			}
		}

	}
}
