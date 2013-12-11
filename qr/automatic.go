package qr

import (
	"fmt"
	"github.com/boombuler/barcode"
)

func encodeAuto(content string, ecl ErrorCorrectionLevel) (*barcode.BitList, *versionInfo, error) {
	bits, vi, _ := Numeric.getEncoder()(content, ecl)
	if bits != nil && vi != nil {
		return bits, vi, nil
	}
	bits, vi, _ = AlphaNumeric.getEncoder()(content, ecl)
	if bits != nil && vi != nil {
		return bits, vi, nil
	}

	return nil, nil, fmt.Errorf("No encoding found to encode \"%s\"", content)
}
