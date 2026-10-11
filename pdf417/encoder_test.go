package pdf417

import "testing"

func TestGetLeftCodeWordClusterZero(t *testing.T) {
	tests := []struct {
		rowNum int
		rows   int
		want   int
	}{
		{rowNum: 0, rows: 4, want: 1},
		{rowNum: 3, rows: 4, want: 31},
		{rowNum: 0, rows: 7, want: 2},
	}

	for _, tt := range tests {
		if got := getLeftCodeWord(tt.rowNum, tt.rows, 3, 0); got != tt.want {
			t.Errorf("getLeftCodeWord(%d, %d, 3, 0) = %d, want %d", tt.rowNum, tt.rows, got, tt.want)
		}
	}
}
