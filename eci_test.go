package barcode

import "testing"

func Test_ECIValid(t *testing.T) {
	tests := map[ECI]bool{
		NoECI - 1:  false,
		NoECI:      true,
		0:          true,
		ECIUTF8:    true,
		MaxECI:     true,
		MaxECI + 1: false,
	}
	for eci, want := range tests {
		if got := eci.Valid(); got != want {
			t.Errorf("ECI(%d).Valid() = %v, want %v", eci, got, want)
		}
	}
}
