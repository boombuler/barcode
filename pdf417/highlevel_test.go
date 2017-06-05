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
	runTest := func(msg string, expected []int) {
		if codes, err := highlevelEncode(msg); err != nil {
			t.Error(err)
		} else {
			compareIntSlice(t, expected, codes, msg)
		}
	}

	runTest("01234", []int{902, 112, 434})
	runTest("Super !", []int{567, 615, 137, 809, 329})
	runTest("Super ", []int{567, 615, 137, 809})
	runTest("ABC123", []int{1, 88, 32, 119})
	runTest("123ABC", []int{841, 63, 840, 32})

	runTest("\x0B", []int{913, 11})
	runTest("\x0B\x0B\x0B\x0B\x0B\x0B", []int{924, 18, 455, 694, 754, 291})
}
