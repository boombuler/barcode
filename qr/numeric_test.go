package qr

import (
	"bytes"
	"strings"
	"testing"
)

func Test_NumericEncoding(t *testing.T) {
	tests := []struct {
		content     string
		wantVersion byte
		wantBytes   []byte
		wantErr     string
	}{
		// Example for Version 1-H symbol from ISO/IEC 18004:2015 §7.4.3
		{
			content:     "01234567",
			wantVersion: 1,
			wantBytes:   []byte{16, 32, 12, 86, 97, 128, 236, 17, 236},
		},
		{
			content:     "0123456789012345",
			wantVersion: 1,
			wantBytes:   []byte{16, 64, 12, 86, 106, 110, 20, 234, 80},
		},
		{content: "foo", wantErr: `"foo" can not be encoded as Numeric`},
		{content: strings.Repeat("1", 14297), wantErr: "To much data to encode"},
	}

	encode := Numeric.getEncoder()
	for _, tt := range tests {
		x, vi, err := encode(tt.content, H)
		if tt.wantErr != "" {
			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("encode(%q) error = %v, want %q", tt.content, err, tt.wantErr)
			}
			if x != nil || vi != nil {
				t.Errorf("encode(%q) returned non-nil result alongside error", tt.content)
			}
			continue
		}
		if err != nil {
			t.Errorf("encode(%q) failed: %v", tt.content, err)
			continue
		}
		if vi.Version != tt.wantVersion {
			t.Errorf("encode(%q) version = %d, want %d", tt.content, vi.Version, tt.wantVersion)
		}
		if got := x.GetBytes(); !bytes.Equal(got, tt.wantBytes) {
			t.Errorf("encode(%q) bytes = %v, want %v", tt.content, got, tt.wantBytes)
		}
	}
}
