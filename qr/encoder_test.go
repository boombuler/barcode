package qr

import (
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/boombuler/barcode"
)

type test struct {
	Text   string
	Mode   Encoding
	ECL    ErrorCorrectionLevel
	Result string
}

var tests = []test{
	{
		Text: "hello world",
		Mode: Unicode,
		ECL:  H,
		Result: `
+++++++.+.+.+...+.+++++++
+.....+.++...+++..+.....+
+.+++.+.+.+.++.++.+.+++.+
+.+++.+....++.++..+.+++.+
+.+++.+..+...++.+.+.+++.+
+.....+.+..+..+++.+.....+
+++++++.+.+.+.+.+.+++++++
........++..+..+.........
..+++.+.+++.+.++++++..+++
+++..+..+...++.+...+..+..
+...+.++++....++.+..++.++
++.+.+.++...+...+.+....++
..+..+++.+.+++++.++++++++
+.+++...+..++..++..+..+..
+.....+..+.+.....+++++.++
+.+++.....+...+.+.+++...+
+.+..+++...++.+.+++++++..
........+....++.+...+.+..
+++++++......++++.+.+.+++
+.....+....+...++...++.+.
+.+++.+.+.+...+++++++++..
+.+++.+.++...++...+.++..+
+.+++.+.++.+++++..++.+..+
+.....+..+++..++.+.++...+
+++++++....+..+.+..+..+++`,
	},
}

func Test_GetUnknownEncoder(t *testing.T) {
	if unknownEncoding.getEncoder() != nil {
		t.Fail()
	}
}

func Test_EncodingStringer(t *testing.T) {
	tests := map[Encoding]string{
		Auto:            "Auto",
		Numeric:         "Numeric",
		AlphaNumeric:    "AlphaNumeric",
		Unicode:         "Unicode",
		unknownEncoding: "",
	}

	for enc, str := range tests {
		if enc.String() != str {
			t.Fail()
		}
	}
}

func Test_InvalidEncoding(t *testing.T) {
	_, err := Encode("hello world", H, Numeric)
	if err == nil {
		t.Fail()
	}
}

func imgStrToBools(str string) []bool {
	res := make([]bool, 0, len(str))
	for _, r := range str {
		switch r {
		case '+':
			res = append(res, true)
		case '.':
			res = append(res, false)
		}
	}
	return res
}

func Test_Encode(t *testing.T) {
	for _, tst := range tests {
		res, err := Encode(tst.Text, tst.ECL, tst.Mode)
		if err != nil {
			t.Error(err)
		}
		qrCode, ok := res.(*qrcode)
		if !ok {
			t.Fail()
		}
		testRes := imgStrToBools(tst.Result)
		if (qrCode.dimension * qrCode.dimension) != len(testRes) {
			t.Fail()
		}
		t.Logf("dim %d", qrCode.dimension)
		for i := 0; i < len(testRes); i++ {
			x := i % qrCode.dimension
			y := i / qrCode.dimension
			if qrCode.Get(x, y) != testRes[i] {
				t.Errorf("Failed at index %d", i)
			}
		}
	}
}

func ExampleEncode() {
	qrcode, err := Encode("hello world", L, Auto)
	if err != nil {
		log.Fatal(err)
	}

	qrcode, err = barcode.Scale(qrcode, 100, 100)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, qrcode); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
