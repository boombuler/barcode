package qr

import (
	"fmt"
	"github.com/boombuler/barcode"
)

type autoEncoding struct {
}

// choose the best matching encoding
var Auto Encoding = autoEncoding{}

func (ne autoEncoding) String() string {
	return "Auto"
}

func (ne autoEncoding) encode(content string, ecl ErrorCorrectionLevel) (*barcode.BitList, *versionInfo, error) {
	bits, vi, _ := Numeric.encode(content, ecl)
	if bits != nil && vi != nil {
		return bits, vi, nil
	}
	bits, vi, _ = AlphaNumeric.encode(content, ecl)
	if bits != nil && vi != nil {
		return bits, vi, nil
	}

	return nil, nil, fmt.Errorf("No encoding found to encode \"%s\"", content)
}
