// Package pdf417 can create PDF-417 barcodes
package pdf417

import (
	"fmt"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

const (
	padding_codeword = 900

	// The codewords that open an ECI, for the designators
	// from 0 to 899, from 900 to 810899, and from 810900 to 811799.
	eci_charset         = 927
	eci_general_purpose = 926
	eci_user_defined    = 925
)

// Encodes the given data and color scheme as PDF417 barcode.
// securityLevel should be between 0 and 8. The higher the number, the more
// additional error-correction codes are added.
func EncodeWithColor(data string, securityLevel byte, color barcode.ColorScheme) (barcode.Barcode, error) {
	return EncodeWithECIAndColor(data, securityLevel, barcode.NoECI, color)
}

// EncodeWithECI encodes the given data as PDF417 barcode like Encode,
// but whose data open with the given ECI, so that a reader knows which
// character set the data's bytes are in.  The data are encoded as the
// bytes they hold, without conversion.  barcode.NoECI writes no ECI.
//
// PDF417 can carry designators from 0 to 811799 only.
func EncodeWithECI(data string, securityLevel byte, eci barcode.ECI) (barcode.Barcode, error) {
	return EncodeWithECIAndColor(data, securityLevel, eci, barcode.ColorScheme16)
}

// EncodeWithECIAndColor encodes the given data as PDF417 barcode
// like EncodeWithECI, in the given color scheme.
func EncodeWithECIAndColor(data string, securityLevel byte, eci barcode.ECI, color barcode.ColorScheme) (barcode.Barcode, error) {
	if securityLevel >= 9 {
		return nil, fmt.Errorf("invalid security level %d", securityLevel)
	}

	if !eci.Valid() {
		return nil, fmt.Errorf("invalid ECI %d", eci)
	}

	sl := securitylevel(securityLevel)

	dataWords, err := dataCodewords(data, eci)
	if err != nil {
		return nil, err
	}

	columns, rows := calcDimensions(len(dataWords), sl.ErrorCorrectionWordCount())
	if columns < minCols || columns > maxCols || rows < minRows || rows > maxRows {
		return nil, fmt.Errorf("unable to fit data in barcode")
	}

	barcode := new(pdfBarcode)
	barcode.data = data
	barcode.color = color

	codeWords, err := encodeData(dataWords, columns, sl)
	if err != nil {
		return nil, err
	}

	grid := [][]int{}
	for i := 0; i < len(codeWords); i += columns {
		grid = append(grid, codeWords[i:min(i+columns, len(codeWords))])
	}

	codes := [][]int{}

	for rowNum, row := range grid {
		table := rowNum % 3
		rowCodes := make([]int, 0, columns+4)

		rowCodes = append(rowCodes, start_word)
		rowCodes = append(rowCodes, getCodeword(table, getLeftCodeWord(rowNum, rows, columns, securityLevel)))

		for _, word := range row {
			rowCodes = append(rowCodes, getCodeword(table, word))
		}

		rowCodes = append(rowCodes, getCodeword(table, getRightCodeWord(rowNum, rows, columns, securityLevel)))
		rowCodes = append(rowCodes, stop_word)

		codes = append(codes, rowCodes)
	}

	barcode.code = renderBarcode(codes)
	barcode.width = (columns+4)*17 + 1

	return barcode, nil
}

// Encodes the given data as PDF417 barcode.
// securityLevel should be between 0 and 8. The higher the number, the more
// additional error-correction codes are added.
func Encode(data string, securityLevel byte) (barcode.Barcode, error) {
	return EncodeWithColor(data, securityLevel, barcode.ColorScheme16)
}

// dataCodewords returns the data codewords of a symbol.
//
// They are the ECI's codewords, if any, followed by
// the high-level encoding of the data.
func dataCodewords(data string, eci barcode.ECI) ([]int, error) {
	eciWords, err := eciCodewords(eci)
	if err != nil {
		return nil, err
	}
	words, err := highlevelEncode(data)
	if err != nil {
		return nil, err
	}
	if len(eciWords) == 0 {
		return words, nil
	}
	return append(eciWords, words...), nil
}

// eciCodewords returns the codewords that open the data of a symbol
// carrying an ECI: the ECI codeword for the designator's range and the
// designator in one or two codewords.  It returns none for barcode.NoECI
// and an error for any other value outside 0 to 811799.
func eciCodewords(eci barcode.ECI) ([]int, error) {
	value := int(eci)
	switch {
	case eci == barcode.NoECI:
		return nil, nil
	case value < 0 || value > 811799:
		return nil, fmt.Errorf(
			"ECI %d is outside 0 to 811799, the range PDF417 can carry", eci)
	case value < 900:
		return []int{eci_charset, value}, nil
	case value < 810900:
		return []int{eci_general_purpose, value/900 - 1, value % 900}, nil
	}
	return []int{eci_user_defined, value - 810900}, nil
}

func encodeData(dataWords []int, columns int, sl securitylevel) ([]int, error) {
	dataCount := len(dataWords)

	ecCount := sl.ErrorCorrectionWordCount()

	padWords := getPadding(dataCount, ecCount, columns)
	dataWords = append(dataWords, padWords...)

	length := len(dataWords) + 1
	dataWords = append([]int{length}, dataWords...)

	ecWords := sl.Compute(dataWords)

	return append(dataWords, ecWords...), nil
}

func getLeftCodeWord(rowNum int, rows int, columns int, securityLevel byte) int {
	tableId := rowNum % 3

	var x int

	switch tableId {
	case 0:
		x = (rows - 1) / 3
	case 1:
		x = int(securityLevel) * 3
		x += (rows - 1) % 3
	case 2:
		x = columns - 1
	}

	return 30*(rowNum/3) + x
}

func getRightCodeWord(rowNum int, rows int, columns int, securityLevel byte) int {
	tableId := rowNum % 3

	var x int

	switch tableId {
	case 0:
		x = columns - 1
	case 1:
		x = (rows - 1) / 3
	case 2:
		x = int(securityLevel) * 3
		x += (rows - 1) % 3
	}

	return 30*(rowNum/3) + x
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func getPadding(dataCount int, ecCount int, columns int) []int {
	totalCount := dataCount + ecCount + 1
	mod := totalCount % columns

	padding := []int{}

	if mod > 0 {
		padCount := columns - mod
		padding = make([]int, padCount)
		for i := 0; i < padCount; i++ {
			padding[i] = padding_codeword
		}
	}

	return padding
}

func renderBarcode(codes [][]int) *utils.BitList {
	bl := new(utils.BitList)
	for _, row := range codes {
		lastIdx := len(row) - 1
		for i, col := range row {
			if i == lastIdx {
				bl.AddBits(col, 18)
			} else {
				bl.AddBits(col, 17)
			}
		}
	}
	return bl
}
