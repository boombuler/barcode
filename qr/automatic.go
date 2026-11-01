package qr

import (
	"fmt"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

func encodeAuto(content string, ecl ErrorCorrectionLevel, eci barcode.ECI) (*utils.BitList, *versionInfo, error) {
	if asciiCompatible(eci) {
		bits, vi, _ := Numeric.getEncoder()(content, ecl, eci)
		if bits != nil && vi != nil {
			return bits, vi, nil
		}
		bits, vi, _ = AlphaNumeric.getEncoder()(content, ecl, eci)
		if bits != nil && vi != nil {
			return bits, vi, nil
		}
	}
	bits, vi, _ := Unicode.getEncoder()(content, ecl, eci)
	if bits != nil && vi != nil {
		return bits, vi, nil
	}
	return nil, nil, fmt.Errorf("no encoding found to encode %q", content)
}

// asciiCompatible reports whether the ECI's character set reads
// the characters of alphanumeric mode as ASCII does.
//
// Numeric and alphanumeric mode write characters, not bytes, so Auto
// may use them only then.  In UTF-16, for example, the bytes "33" are
// one character, U+3333, not two digits.  A designator not known to be
// such a character set gets byte mode, which writes the bytes as they are.
func asciiCompatible(eci barcode.ECI) bool {
	return eci == barcode.NoECI ||
		eci >= 0 && eci <= barcode.ECIISO8859_11 ||
		eci >= barcode.ECIISO8859_13 && eci <= barcode.ECIISO8859_16 ||
		eci >= barcode.ECIShiftJIS && eci <= barcode.ECIWindows1256 ||
		eci >= barcode.ECIUTF8 && eci <= barcode.ECIGB18030
}
