// Package pdf417 can create PDF-417 barcodes
package pdf417

import (
	"fmt"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/utils"
)

const (
	padding_codeword = 900
)

// Encodes the given data and color scheme as PDF417 barcode.
// securityLevel should be between 0 and 8. The higher the number, the more
// additional error-correction codes are added.
func EncodeWithColor(data string, securityLevel byte, color barcode.ColorScheme) (barcode.Barcode, error) {
	if securityLevel >= 9 {
		return nil, fmt.Errorf("invalid security level %d", securityLevel)
	}

	sl := securitylevel(securityLevel)

	dataWords, err := highlevelEncode(data)
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

	codeWords := encodeData(dataWords, columns, rows, sl)

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

// encodeData returns the codewords of the PDF417 symbol matrix ready
// for low-level encodation: padded data codewords for source data,
// followed by error correction codewords.
func encodeData(sourceData []int, columns int, rows int, sl securitylevel) []int {
	// number of error correction codewords
	k := sl.ErrorCorrectionWordCount()
	// number of data codewords excluding padding and symbol length descriptor
	m := len(sourceData)
	// number of data codewords including padding and symbol length descriptor
	// (c * r) = (n + k), so n = (c * r) - k
	n := columns*rows - k

	// Allocate data including error correction:
	// <Symbol Length Descriptor> <Data Codewords> [Pad Codewords] <EC Codewords>
	data := make([]int, 0, n+k)
	data = append(data, n)             // Symbol Length Descriptor
	data = append(data, sourceData...) // Data Codewords

	// §5.9.2: pad codewords = (n - m) - 1
	padCount := (n - m) - 1
	if padCount < 0 {
		panic(fmt.Sprintf("pdf417: %d data codewords do not fit %d columns, %d rows at level %d", m, columns, rows, sl))
	}
	for i := 0; i < padCount; i++ {
		data = append(data, padding_codeword)
	}

	// §5.10: error correction codewords calculated from Symbol Length Descriptor, data, and pad.
	return append(data, sl.Compute(data)...)
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
