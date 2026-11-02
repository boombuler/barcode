package pdf417

import (
	"fmt"
	"strings"
	"testing"
)

func TestRowIndicatorCodewords(t *testing.T) {
	tests := []struct {
		name          string
		rowNum        int
		rows          int
		columns       int
		securityLevel byte
		wantLeft      int
		wantRight     int
	}{
		{"table0/group0", 0, 17, 6, 3, 5, 5},
		{"table1/group0", 1, 17, 6, 3, 10, 5},
		{"table2/group0", 2, 17, 6, 3, 5, 10},
		{"table0/group1", 3, 17, 6, 3, 35, 35},
		{"table1/group1", 4, 17, 6, 3, 40, 35},
		{"table2/group1", 5, 17, 6, 3, 35, 40},
		{"table1/group5", 16, 17, 6, 3, 160, 155},
		{"table0/rowsMod1", 0, 16, 6, 3, 5, 5},
		{"table0/maxParams", 0, maxRows, maxCols, 8, 9, 29},
		{"table1/maxParams", 1, maxRows, maxCols, 8, 26, 9},
		{"table2/maxParams", 2, maxRows, maxCols, 8, 29, 26},
		{"table0/lastRow", maxRows - 3, maxRows, maxCols, 8, 279, 299},
		{"table1/lastRow", maxRows - 2, maxRows, maxCols, 8, 296, 279},
		{"table2/lastRow", maxRows - 1, maxRows, maxCols, 8, 299, 296},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLeftCodeWord(tt.rowNum, tt.rows, tt.columns, tt.securityLevel); got != tt.wantLeft {
				t.Errorf("getLeftCodeWord(rowNum=%d, rows=%d, columns=%d, securityLevel=%d) = %d, want %d",
					tt.rowNum, tt.rows, tt.columns, tt.securityLevel, got, tt.wantLeft)
			}
			if got := getRightCodeWord(tt.rowNum, tt.rows, tt.columns, tt.securityLevel); got != tt.wantRight {
				t.Errorf("getRightCodeWord(rowNum=%d, rows=%d, columns=%d, securityLevel=%d) = %d, want %d",
					tt.rowNum, tt.rows, tt.columns, tt.securityLevel, got, tt.wantRight)
			}
		})
	}
}

func TestMinimumRows(t *testing.T) {
	// ISO/IEC 15438 symbols have at least 3 rows; readers reject 2-row symbols
	for _, data := range []string{"0", "A", "AB", "Hell", "0123"} {
		for securityLevel := byte(0); securityLevel <= 2; securityLevel++ {
			code, err := Encode(data, securityLevel)
			if err != nil {
				t.Fatalf("Encode(%q, %d): %v", data, securityLevel, err)
			}
			if rows := code.Bounds().Dy() / moduleHeight; rows < 3 {
				t.Errorf("Encode(%q, %d) has %d rows, want at least 3", data, securityLevel, rows)
			}
		}
	}
}

func TestEncodeDataFillsSymbol(t *testing.T) {
	// The codewords fill the symbol: the symbol length descriptor, which
	// counts itself, the data and the pad codewords, then the data, then pad
	// codewords up to the error correction codewords at the end.
	for _, data := range []string{"A", "AB", "Hell", "0123", "Hello, World!", strings.Repeat("1234567890", 10)} {
		for securityLevel := byte(0); securityLevel <= 4; securityLevel++ {
			sl := securitylevel(securityLevel)
			dataWords, err := highlevelEncode(data)
			if err != nil {
				t.Fatalf("highlevelEncode(%q): %v", data, err)
			}
			ecCount := sl.ErrorCorrectionWordCount()
			columns, rows := calcDimensions(len(dataWords), ecCount)

			codeWords, err := encodeData(dataWords, columns, rows, sl)
			if err != nil {
				t.Fatalf("encodeData(%q, %d): %v", data, securityLevel, err)
			}

			name := fmt.Sprintf("%q at level %d (%d columns, %d rows)", data, securityLevel, columns, rows)
			if len(codeWords) != columns*rows {
				t.Errorf("%s: %d codewords, want %d", name, len(codeWords), columns*rows)
				continue
			}
			if want := columns*rows - ecCount; codeWords[0] != want {
				t.Errorf("%s: symbol length descriptor %d, want %d", name, codeWords[0], want)
			}
			for i, word := range dataWords {
				if codeWords[1+i] != word {
					t.Errorf("%s: codeword %d is %d, want the data codeword %d", name, 1+i, codeWords[1+i], word)
				}
			}
			for i, word := range codeWords[1+len(dataWords) : columns*rows-ecCount] {
				if word != padding_codeword {
					t.Errorf("%s: codeword %d is %d, want the pad codeword %d", name, 1+len(dataWords)+i, word, padding_codeword)
				}
			}
		}
	}
}
