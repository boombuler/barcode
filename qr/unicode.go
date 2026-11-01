package qr

import (
	"errors"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

func encodeUnicode(content string, ecl ErrorCorrectionLevel, eci barcode.ECI) (*utils.BitList, *versionInfo, error) {
	data := []byte(content)

	// Without an ECI header a reader that follows the standard reads
	// the bytes as ISO 8859-1, not UTF-8.  Some readers cannot handle
	// the header, so it is written only when the caller asks for one.
	res := eciHeader(eci)
	vi := findSmallestVersionInfo(ecl, byteMode, res.Len()+len(data)*8)
	if vi == nil {
		return nil, nil, errors.New("too much data to encode")
	}

	res.AddBits(int(byteMode), 4)
	res.AddBits(len(content), vi.charCountBits(byteMode))
	for _, b := range data {
		res.AddByte(b)
	}
	addPaddingAndTerminator(res, vi)
	return res, vi, nil
}
