package pdf417

import "testing"

func TestEncode(t *testing.T) {
	encoder := newDataEncoder()

	// When starting with text, the first code word does not need to be the switch
	if result, err := encoder.Encode("ABC123"); err != nil {
		t.Error(err)
	} else {
		compareIntSlice(t, []int{1, 89, 902, 1, 223}, result)
	}
	// When starting with numbers, we do need to switchresult := encoder.Encode("ABC123")
	if result, err := encoder.Encode("123ABC"); err != nil {
		t.Error(err)
	} else {
		compareIntSlice(t, []int{902, 1, 223, 900, 1, 89}, result)
	}

	/*
		// Also with bytes
		if result, err := encoder.Encode("\x0B"); err != nil {
			t.Error(err)
		} else {
			compareIntSlice(t, []int{901, 11}, result)
		}

		// Alternate bytes switch code when number of bytes is divisble by 6
		if result, err := encoder.Encode("\x0B\x0B\x0B\x0B\x0B\x0B"); err != nil {
			t.Error(err)
		} else {
			compareIntSlice(t, []int{924, 18, 455, 694, 754, 291}, result)
		}
	*/
}
