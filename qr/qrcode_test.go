package qr

import (
	"image/color"
	"testing"
)

func Test_NewQRCode(t *testing.T) {
	bc := newBarcode(2)
	if bc == nil {
		t.Fail()
	}
	if bc.data.Len() != 4 {
		t.Fail()
	}
	if bc.dimension != 2 {
		t.Fail()
	}
}

func Test_ImageBasics(t *testing.T) {
	qr := newBarcode(10)
	if qr.ColorModel() != color.Gray16Model {
		t.Fail()
	}
}



func Test_Penalty4(t *testing.T) {
	qr := newBarcode(3)
	if qr.calcPenaltyRule4() != 100 {
		t.Fail()
	}
	qr.Set(0, 0, true)
	if qr.calcPenaltyRule4() != 70 {
		t.Fail()
	}
	qr.Set(0, 1, true)
	if qr.calcPenaltyRule4() != 50 {
		t.Fail()
	}
	qr.Set(0, 2, true)
	if qr.calcPenaltyRule4() != 30 {
		t.Fail()
	}
	qr.Set(1, 0, true)
	if qr.calcPenaltyRule4() != 10 {
		t.Fail()
	}
	qr.Set(1, 1, true)
	if qr.calcPenaltyRule4() != 10 {
		t.Fail()
	}
	qr = newBarcode(2)
	qr.Set(0, 0, true)
	qr.Set(1, 0, true)
	if qr.calcPenaltyRule4() != 0 {
		t.Fail()
	}
}
