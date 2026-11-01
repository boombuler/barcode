// Package datamatrix can create Datamatrix barcodes
package datamatrix

import (
	"errors"
	"fmt"

	"github.com/boombuler/barcode"
)

// FNC1 is the codeword for the Function 1 Symbol Character to
// differentiate a GS1 DataMatrix from other Data Matrix symbols.
//
// It is used as both a start character and a separator of GS1 element
// strings.
const FNC1 byte = 232

// eciCodeword is the codeword that opens an ECI.
const eciCodeword byte = 241

// Encode returns a Datamatrix barcode for the given content and color scheme
func EncodeWithColor(content string, color barcode.ColorScheme) (barcode.Barcode, error) {
	return EncodeWithECIAndColor(content, barcode.NoECI, color)
}

// EncodeWithECI returns a Datamatrix barcode like Encode, but whose data
// open with the given ECI, so that a reader knows which character set
// the content's bytes are in.  The content is encoded as the bytes it
// holds, without conversion.  barcode.NoECI writes no ECI.
//
// A GS1 symbol cannot carry an ECI, so the content is never read as GS1:
// a leading FNC1 byte is encoded as data, like every other byte.
// Use Encode for GS1 content.  This way UTF-8 text can start with a
// character from U+8000 to U+8FFF, such as many common CJK ideographs,
// whose first byte is the same as FNC1.
func EncodeWithECI(content string, eci barcode.ECI) (barcode.Barcode, error) {
	return EncodeWithECIAndColor(content, eci, barcode.ColorScheme16)
}

// EncodeWithECIAndColor returns a Datamatrix barcode like EncodeWithECI,
// in the given color scheme.
func EncodeWithECIAndColor(content string, eci barcode.ECI, color barcode.ColorScheme) (barcode.Barcode, error) {
	if !eci.Valid() {
		return nil, fmt.Errorf("invalid ECI %d", eci)
	}
	data := encodeText(content, eci)

	var size *dmCodeSize
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(data) {
			size = s
			break
		}
	}
	if size == nil {
		return nil, errors.New("too much data to encode")
	}
	data = addPadding(data, size.DataCodewords())
	data = ec.calcECC(data, size)
	code := render(data, size, color)
	if code != nil {
		code.content = content
		return code, nil
	}
	return nil, errors.New("unable to render barcode")
}

// Encode returns a Datamatrix barcode for the given content
func Encode(content string) (barcode.Barcode, error) {
	return EncodeWithColor(content, barcode.ColorScheme16)
}

func render(data []byte, size *dmCodeSize, color barcode.ColorScheme) *datamatrixCode {
	cl := newCodeLayout(size, color)

	cl.SetValues(data)

	return cl.Merge()
}

func encodeText(content string, eci barcode.ECI) []byte {
	result := eciCodewords(eci)
	input := []byte(content)

	isGS1 := eci == barcode.NoECI && len(input) > 0 && input[0] == FNC1

	for i := 0; i < len(input); {
		c := input[i]
		i++

		if c >= '0' && c <= '9' && i < len(input) && input[i] >= '0' && input[i] <= '9' {
			// two numbers...
			c2 := input[i]
			i++
			cw := byte(((c-'0')*10 + (c2 - '0')) + 130)
			result = append(result, cw)
		} else if isGS1 && c == FNC1 {
			result = append(result, c)
		} else if c > 127 {
			// not correct... needs to be redone later...
			result = append(result, 235, c-127)
		} else {
			result = append(result, c+1)
		}
	}
	return result
}

// eciCodewords returns the codewords that open the data of a symbol
// carrying an ECI: the ECI codeword and the designator in one, two,
// or three codewords.  It returns none for barcode.NoECI.
func eciCodewords(eci barcode.ECI) []byte {
	value := int(eci)
	switch {
	case eci == barcode.NoECI:
		return nil
	case value <= 126:
		return []byte{eciCodeword, byte(value + 1)}
	case value <= 16382:
		value -= 127
		return []byte{eciCodeword, byte(value/254 + 128), byte(value%254 + 1)}
	default:
		value -= 16383
		return []byte{eciCodeword, byte(value/64516 + 192),
			byte(value/254%254 + 1), byte(value%254 + 1)}
	}
}

func addPadding(data []byte, toCount int) []byte {
	if len(data) < toCount {
		data = append(data, 129)
	}
	for len(data) < toCount {
		R := ((149 * (len(data) + 1)) % 253) + 1
		tmp := 129 + R
		if tmp > 254 {
			tmp = tmp - 254
		}

		data = append(data, byte(tmp))
	}
	return data
}
