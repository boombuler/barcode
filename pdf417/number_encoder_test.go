package pdf417

import (
	"testing"
)

func TestNumberEncoder_CanEncode(t *testing.T) {
	encoder := newNumberEncoder()

	shouldEncode := func(tests ...rune) {
		for _, test := range tests {
			if !encoder.CanEncode(test) {
				t.Errorf("NumberEncoder should be able to encode %q", string(test))
			}
		}
	}
	shouldEncode('0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	shouldNotEncode := func(tests ...rune) {
		for _, test := range tests {
			if encoder.CanEncode(test) {
				t.Errorf("NumberEncoder should not be able to encode %q", string(test))
			}
		}
	}
	shouldNotEncode('a', 'q', '\t')
}

func TestNumberEncoder_GetSwitchCode(t *testing.T) {
	encoder := newNumberEncoder()
	if sc := encoder.GetSwitchCode("123"); sc != switchCodeNumber {
		t.Errorf("Unexpected switchcode. Got %v", sc)
	}
	if sc := encoder.GetSwitchCode("foo"); sc != switchCodeNumber {
		t.Errorf("Unexpected switchcode. Got %v", sc)
	}
}

func TestNumberEncoder_Encode(t *testing.T) {
	encoder := newNumberEncoder()

	if codes, err := encoder.Encode("01234"); err != nil {
		t.Error(err)
	} else {
		compareIntSlice(t, []int{112, 434}, codes)
	}
}
