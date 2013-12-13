package qr

import (
	"bytes"
	"testing"
)

func Test_LogTables(t *testing.T) {
	for i := 1; i <= 255; i++ {
		tmp := ec.fld.LogTbl[i]
		if i != ec.fld.ALogTbl[tmp] {
			t.Errorf("Invalid LogTables: %d", i)
		}
	}

	if ec.fld.ALogTbl[11] != 232 || ec.fld.ALogTbl[87] != 127 || ec.fld.ALogTbl[225] != 36 {
		t.Fail()
	}
}

func Test_GetPolynomial(t *testing.T) {
	doTest := func(b []byte) {
		cnt := byte(len(b) - 1)
		if bytes.Compare(ec.getPolynomial(cnt), b) != 0 {
			t.Errorf("Failed getPolynomial(%d)", cnt)
		}
	}
	doTest([]byte{0, 0})
	doTest([]byte{0, 87, 229, 146, 149, 238, 102, 21})
	doTest([]byte{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45})
	doTest([]byte{0, 183, 26, 201, 87, 210, 221, 113, 21, 46, 65, 45, 50, 238, 184, 249, 225, 102, 58, 209, 218, 109, 165, 26, 95, 184, 192, 52, 245, 35, 254, 238, 175, 172, 79, 123, 25, 122, 43, 120, 108, 215, 80, 128, 201, 235, 8, 153, 59, 101, 31, 198, 76, 31, 156})
}

func Test_ErrorCorrection(t *testing.T) {
	doTest := func(b []byte, ecc []byte) {
		cnt := byte(len(ecc))
		if bytes.Compare(ec.calcECC(b, cnt), ecc) != 0 {
			t.Fail()
		}
	}

	data1 := []byte{16, 32, 12, 86, 97, 128, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17}
	doTest(data1, []byte{140, 250})
	doTest(data1, []byte{165, 36, 212, 193, 237, 54, 199, 135, 44, 85})
	doTest(data1, []byte{227, 219, 167, 206, 127, 77, 181, 205, 203, 131, 6, 102, 62, 113, 173, 153, 69, 210, 55, 111, 146, 227, 13, 144, 249, 87, 103, 81, 30, 125, 189, 61, 142, 129, 129, 43, 148, 88, 25, 249, 37, 58, 57, 108, 91, 241, 78, 248, 226, 177, 17, 58, 59, 218, 53, 146, 96, 165, 146, 163, 198, 190, 15, 71, 117, 164, 167, 53})

	data2 := []byte{0, 0, 255, 255}
	doTest(data2, []byte{171, 81, 216, 241, 210})
	doTest(data2, []byte{12, 183, 205, 34, 73, 117, 36, 75, 237, 235})
}
