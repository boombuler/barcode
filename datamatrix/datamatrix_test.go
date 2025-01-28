package datamatrix

import (
	"bytes"
	"testing"
)

func codeFromStr(str string, size *dmCodeSize) *datamatrixCode {
	code := newDataMatrixCode(size)
	idx := 0
	for _, r := range str {
		x := idx % size.Columns
		y := idx / size.Columns

		switch r {
		case '#':
			code.set(x, y, true)
		case '.':
			code.set(x, y, false)
		default:
			continue
		}

		idx++
	}
	return code
}

func Test_Issue12(t *testing.T) {
	data := `{"po":12,"batchAction":"start_end"}`
	realData := addPadding(encodeText(data), 36)
	wantedData := []byte{124, 35, 113, 112, 35, 59, 142, 45, 35, 99, 98, 117, 100, 105, 66, 100, 117, 106, 112, 111, 35, 59, 35, 116, 117, 98, 115, 117, 96, 102, 111, 101, 35, 126, 129, 181}

	if bytes.Compare(realData, wantedData) != 0 {
		t.Error("Data Encoding failed")
		return
	}

	var codeSize *dmCodeSize
	for _, s := range codeSizes {
		if s.DataCodewords() >= len(wantedData) {
			codeSize = s
			break
		}
	}
	realECC := ec.calcECC(realData, codeSize)[len(realData):]
	wantedECC := []byte{196, 53, 147, 192, 151, 213, 107, 61, 98, 251, 50, 71, 186, 15, 43, 111, 165, 243, 209, 79, 128, 109, 251, 4}
	if bytes.Compare(realECC, wantedECC) != 0 {
		t.Errorf("Error correction calculation failed\nGot: %v", realECC)
		return
	}

	barcode := `
#.#.#.#.#.#.#.#.#.#.#.#.
#....###..#..#....#...##
##.......#...#.#.#....#.
#.###...##..#...##.##..#
##...####..##..#.#.#.##.
#.###.##.###..#######.##
#..###...##.##..#.##.##.
#.#.#.#.#.#.###....#.#.#
##.#...#.#.#..#...#####.
#...####..#...##..#.#..#
##...#...##.###.#.....#.
#.###.#.##.#.....###..##
##..#####...#..##...###.
###...#.####.##.#.#.#..#
#..###..#.#.####.#.###..
###.#.#..#..#.###.#.##.#
#####.##.###..#.####.#..
#.##.#......#.#..#.#.###
###.#....######.#...##..
##...#..##.###..#...####
#.######.###.##..#...##.
#..#..#.##.#..####...#.#
###.###..#..##.#.##...#.
########################`

	bc, err := Encode(data)

	if err != nil {
		t.Error(err)
		return
	}
	realResult := bc.(*datamatrixCode)
	if realResult.Columns != 24 || realResult.Rows != 24 {
		t.Errorf("Got wrong barcode size %dx%d", realResult.Columns, realResult.Rows)
		return
	}

	wantedResult := codeFromStr(barcode, realResult.dmCodeSize)

	for x := 0; x < wantedResult.Columns; x++ {
		for y := 0; y < wantedResult.Rows; y++ {
			r := realResult.get(x, y)
			w := wantedResult.get(x, y)
			if w != r {
				t.Errorf("Failed at: c%d/r%d", x, y)
			}
		}
	}
}

func Test_GS1DataMatrix(t *testing.T) {
	// Example 2 from the GS1 DataMatrix Guideline.
	//
	// (01)09501101020917(17)190508(10)ABCD1234(21)10
	//
	// See: https://www.gs1.org/standards/gs1-datamatrix-guideline/25#2-Encoding-data+2-3-Human-readable-interpretation-(HRI)
	data := new(bytes.Buffer)
	data.WriteByte(FNC1)                 // Start Character
	data.WriteString("0109501101020917") // AI (01)
	data.WriteString("17190508")         // AI (17)
	data.WriteString("10ABCD1234")       // AI (10) does not have pre-defined length
	data.WriteByte(FNC1)                 // Separator Character
	data.WriteString("2110")             // AI (20)

	// Codewords from decoding example 2 with "dmtxread -c".
	wantedData := []byte{
		232, // FNC1
		131, 139, 180, 141, 131, 132, 139, 147,
		147, 149, 135, 138,
		140, 66, 67, 68, 69, 142, 164,
		232, // FNC1
		151, 140,
	}

	realData := encodeText(data.String())
	if bytes.Compare(realData, wantedData) != 0 {
		t.Errorf("GS1 DataMatrix encoding failed\nwant: %v\ngot:  %v\n", wantedData, realData)
	}
}
