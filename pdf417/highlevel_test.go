package pdf417

import "testing"

func compareIntSlice(t *testing.T, expected, actual []int, testStr string) {
	if len(actual) != len(expected) {
		t.Errorf("Invalid slice size. Expected %d got %d while encoding %q", len(expected), len(actual), testStr)
		return
	}
	for i, a := range actual {
		if e := expected[i]; e != a {
			t.Errorf("Unexpected value at position %d. Expected %d got %d while encoding %q", i, e, a, testStr)
		}
	}
}

func TestHighlevelEncode(t *testing.T) {
	runTest := func(msg string, expected ...int) {
		if codes, err := highlevelEncode(msg); err != nil {
			t.Error(err)
		} else {
			compareIntSlice(t, expected, codes, msg)
		}
	}

	runTest("01234", 902, 112, 434)
	runTest("Super !", 567, 615, 137, 809, 329)
	runTest("Super ", 567, 615, 137, 809)
	runTest("ABC123", 1, 88, 32, 119)
	runTest("123ABC", 841, 63, 840, 32)

	// Add Text Compaction examples from ISO/IEC 15438:2015 Table 6.
	runTest("PD", 453)
	runTest("F41", 178, 121)    // combined "F<ml>" and "41" examples
	runTest("A7;", 28, 239, 29) // "A" to latch mixed, then ";" to shift punc
	runTest("PDF417", 453, 178, 121, 239)
}

func TestBinaryEncoder(t *testing.T) {
	runTest := func(msg string, expected ...int) {
		codes := encodeBinary([]byte(msg), encText)
		compareIntSlice(t, expected, codes, msg)
	}

	runTest("alcool", 924, 163, 238, 432, 766, 244)
	runTest("alcoolique", 901, 163, 238, 432, 766, 244, 105, 113, 117, 101)
}
