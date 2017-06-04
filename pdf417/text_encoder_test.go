package pdf417

import (
	"testing"
)

func compareIntSlice(t *testing.T, expected, actual []int) {
	if len(actual) != len(expected) {
		t.Errorf("Invalid slice size. Expected %d got %d", len(expected), len(actual))
		return
	}
	for i, a := range actual {
		if e := expected[i]; e != a {
			t.Errorf("Unexpected value at position %d. Expected %d got %d", i, e, a)
		}
	}
}

func TestTextEncoder_CanEncode(t *testing.T) {
	encoder := newTextEncoder()

	for ord := int(' '); ord < int('Z'); ord++ {
		chr := rune(ord)

		if chr == '"' {
			continue
		}

		if !encoder.CanEncode(chr) {
			t.Errorf("Unable to encode: %d %c", ord, chr)
		}
	}
}

func TestTextEncoder_GetSwitchCode(t *testing.T) {
	encoder := newTextEncoder()
	if sc := encoder.GetSwitchCode("123"); sc != switchCodeText {
		t.Errorf("Unexpected switchcode. Got %v", sc)
	}
	if sc := encoder.GetSwitchCode("foo"); sc != switchCodeText {
		t.Errorf("Unexpected switchcode. Got %v", sc)
	}
}

func TestTextEncoder_Encode(t *testing.T) {
	encoder := newTextEncoder()

	if codes, err := encoder.Encode("Super !"); err != nil {
		t.Error(err)
	} else {
		compareIntSlice(t, []int{567, 615, 137, 808, 760}, codes)
	}

	if codes, err := encoder.Encode("Super "); err != nil {
		t.Error(err)
	} else {
		compareIntSlice(t, []int{567, 615, 137, 809}, codes)
	}
}
