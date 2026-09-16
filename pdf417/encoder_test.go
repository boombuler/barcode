package pdf417

import "testing"

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
