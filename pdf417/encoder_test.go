package pdf417

import (
	"reflect"
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

// TestEncodeData checks the data codewords that §5.9.2 lays out, and that
// error correction codewords complete the symbol (§5.10). The columns,
// rows and level are inputs, so each case chooses the layout it tests.
func TestEncodeData(t *testing.T) {
	tests := []struct {
		name          string
		data          []int
		columns, rows int
		level         securitylevel
		// Symbol Length Descriptor, the source data, any pad codewords
		want []int
	}{
		{
			// The smallest symbol, 2x3, holds up to 3 data words at level 0.
			name:    "no pad",
			data:    []int{10, 20, 30},
			columns: 2,
			rows:    3,
			level:   0,
			want:    []int{4, 10, 20, 30},
		},
		{
			// The smallest symbol, 2x3, has 1 pad with 2 data words.
			name:    "one pad",
			data:    []int{10, 20},
			columns: 2,
			rows:    3,
			level:   0,
			want:    []int{4, 10, 20, 900},
		},
		{
			// The smallest symbol, 2x3, has 2 pads with 1 data word.
			name:    "two pads",
			data:    []int{10},
			columns: 2,
			rows:    3,
			level:   0,
			want:    []int{4, 10, 900, 900},
		},
		{
			// The pad value 900 is still data, so it should be followed by a pad.
			// This ensures the test measures padding correctly.
			name:    "data word with the pad value",
			data:    []int{10, 900},
			columns: 2,
			rows:    3,
			level:   0,
			want:    []int{4, 10, 900, 900},
		},
		{
			// A 3x3 symbol holds 9 codewords, with level 1 using 4 for error
			// correction, so it can fit 4 data words without padding.
			name:    "level 1, no pad",
			data:    []int{10, 20, 30, 40},
			columns: 3,
			rows:    3,
			level:   1,
			want:    []int{5, 10, 20, 30, 40},
		},
		{
			// A 3x3 symbol holds 9 codewords, with level 1 using 4 for error
			// correction, so 3 data words leave 1 pad.
			name:    "level 1, one pad",
			data:    []int{10, 20, 30},
			columns: 3,
			rows:    3,
			level:   1,
			want:    []int{5, 10, 20, 30, 900},
		},
		{
			// The highest level uses 512 error correction codewords,
			// so 20x26 = 520 leaves n = 8: the descriptor, 1 data word and 6 pads.
			name:    "level 8",
			data:    []int{10},
			columns: 20,
			rows:    26,
			level:   8,
			want:    []int{8, 10, 900, 900, 900, 900, 900, 900},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeData(tt.data, tt.columns, tt.rows, tt.level)
			if len(got) != tt.columns*tt.rows {
				t.Fatalf("%d codewords, want %d to fill %dx%d", len(got), tt.columns*tt.rows, tt.columns, tt.rows)
			}

			k := tt.level.ErrorCorrectionWordCount()
			if data := got[:len(got)-k]; !reflect.DeepEqual(data, tt.want) {
				t.Errorf("data codewords (n = %d, m = %d) = %v, want %v", len(tt.want), len(tt.data), data, tt.want)
			}
			checkErrorCorrection(t, got, k)
		})
	}
}

// checkErrorCorrection checks the error correction codewords against
// §5.10. The codewords are the coefficients of a polynomial C(x), highest
// power first, which must be divisible by the generator polynomial
//
//	g(x) = (x - 3)(x - 3^2)...(x - 3^k)
//
// so C(x) is zero at each root 3^i when there are no errors.
func checkErrorCorrection(t *testing.T, codewords []int, k int) {
	t.Helper()

	const modulus = 929 // §5.10: GF(929)

	root := 1
	for i := 1; i <= k; i++ {
		root = root * 3 % modulus // 3^i

		// Horner's method: C(x) = (...((c0*x + c1)*x + c2)...)*x + cN
		got := 0
		for _, c := range codewords {
			got = (got*root + c) % modulus
		}
		if got != 0 {
			t.Errorf("C(3^%d) = %d, want 0", i, got)
		}
	}
}
