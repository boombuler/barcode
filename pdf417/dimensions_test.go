package pdf417

import "testing"

// TestDimensionsFitData checks that every size calcDimensions chooses
// within the symbol limits has room for the data. Only the data
// length and level affect the size, so every combination is checked.
func TestDimensionsFitData(t *testing.T) {
	const (
		// Annex O: (n + k) = (c * r) < 929
		maxCodewords = 928
		// At level 0: less the descriptor and 2 error correction codewords
		maxDataCodewords = maxCodewords - 1 - 2
	)
	for level := securitylevel(0); level <= 8; level++ {
		k := level.ErrorCorrectionWordCount()
		for m := 0; m <= maxDataCodewords; m++ {
			c, r := calcDimensions(m, k)
			if c < minCols || c > maxCols || r < minRows || r > maxRows {
				continue // EncodeWithColor rejects this size
			}
			if c*r > maxCodewords {
				t.Errorf("calcDimensions(%d, %d) = %dx%d, which holds %d codewords, want at most %d", m, k, c, r, c*r, maxCodewords)
			}
			if c*r < m+1+k {
				t.Errorf("calcDimensions(%d, %d) = %dx%d, which holds %d codewords, want at least %d", m, k, c, r, c*r, m+1+k)
			}
		}
	}
}
