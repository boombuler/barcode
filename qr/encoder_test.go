package qr

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/boombuler/barcode"
)

var qrHelloWorldHUni = []bool{true, true, true, true, true, true, true, false, true, false, true, false, true, false, false, false, true, false, true, true, true, true, true, true, true,
	true, false, false, false, false, false, true, false, true, true, false, false, false, true, true, true, false, false, true, false, false, false, false, false, true,
	true, false, true, true, true, false, true, false, true, false, true, false, true, true, false, true, true, false, true, false, true, true, true, false, true,
	true, false, true, true, true, false, true, false, false, false, false, true, true, false, true, true, false, false, true, false, true, true, true, false, true,
	true, false, true, true, true, false, true, false, false, true, false, false, false, true, true, false, true, false, true, false, true, true, true, false, true,
	true, false, false, false, false, false, true, false, true, false, false, true, false, false, true, true, true, false, true, false, false, false, false, false, true,
	true, true, true, true, true, true, true, false, true, false, true, false, true, false, true, false, true, false, true, true, true, true, true, true, true,
	false, false, false, false, false, false, false, false, true, true, false, false, true, false, false, true, false, false, false, false, false, false, false, false, false,
	false, false, true, true, true, false, true, false, true, true, true, false, true, false, true, true, true, true, true, true, false, false, true, true, true,
	true, true, true, false, false, true, false, false, true, false, false, false, true, true, false, true, false, false, false, true, false, false, true, false, false,
	true, false, false, false, true, false, true, true, true, true, false, false, false, false, true, true, false, true, false, false, true, true, false, true, true,
	true, true, false, true, false, true, false, true, true, false, false, false, true, false, false, false, true, false, true, false, false, false, false, true, true,
	false, false, true, false, false, true, true, true, false, true, false, true, true, true, true, true, false, true, true, true, true, true, true, true, true,
	true, false, true, true, true, false, false, false, true, false, false, true, true, false, false, true, true, false, false, true, false, false, true, false, false,
	true, false, false, false, false, false, true, false, false, true, false, true, false, false, false, false, false, true, true, true, true, true, false, true, true,
	true, false, true, true, true, false, false, false, false, false, true, false, false, false, true, false, true, false, true, true, true, false, false, false, true,
	true, false, true, false, false, true, true, true, false, false, false, true, true, false, true, false, true, true, true, true, true, true, true, false, false,
	false, false, false, false, false, false, false, false, true, false, false, false, false, true, true, false, true, false, false, false, true, false, true, false, false,
	true, true, true, true, true, true, true, false, false, false, false, false, false, true, true, true, true, false, true, false, true, false, true, true, true,
	true, false, false, false, false, false, true, false, false, false, false, true, false, false, false, true, true, false, false, false, true, true, false, true, false,
	true, false, true, true, true, false, true, false, true, false, true, false, false, false, true, true, true, true, true, true, true, true, true, false, false,
	true, false, true, true, true, false, true, false, true, true, false, false, false, true, true, false, false, false, true, false, true, true, false, false, true,
	true, false, true, true, true, false, true, false, true, true, false, true, true, true, true, true, false, false, true, true, false, true, false, false, true,
	true, false, false, false, false, false, true, false, false, true, true, true, false, false, true, true, false, true, false, true, true, false, false, false, true,
	true, true, true, true, true, true, true, false, false, false, false, true, false, false, true, false, true, false, false, true, false, false, true, true, true,
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

func Test_Encode(t *testing.T) {
	res, err := Encode("hello world", H, Unicode)
	if err != nil {
		t.Error(err)
	}
	qrCode, ok := res.(*qrcode)
	if !ok {
		t.Fail()
	}
	if (qrCode.dimension * qrCode.dimension) != len(qrHelloWorldHUni) {
		t.Fail()
	}
	t.Logf("dim %d", qrCode.dimension)
	for i := 0; i < len(qrHelloWorldHUni); i++ {
		x := i % qrCode.dimension
		y := i / qrCode.dimension
		if qrCode.Get(x, y) != qrHelloWorldHUni[i] {
			t.Errorf("Failed at index %d", i)
		}
	}
}

func ExampleEncode() {
	f, _ := os.Create("qrcode.png")
	defer f.Close()

	qrcode, err := Encode("hello world", L, Auto)
	if err != nil {
		fmt.Println(err)
	} else {
		qrcode, err = barcode.Scale(qrcode, 100, 100)
		if err != nil {
			fmt.Println(err)
		} else {
			png.Encode(f, qrcode)
		}
	}
}
