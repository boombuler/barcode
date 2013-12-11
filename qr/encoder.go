package qr

import (
	"fmt"
	"github.com/boombuler/barcode"
	"image"
)

type encodeFn func(content string, eccLevel ErrorCorrectionLevel) (*barcode.BitList, *versionInfo, error)

type Encoding byte

const (
	// Choose best matching encoding
	Auto Encoding = iota
	// Encode only numbers [0-9]
	Numeric
	// Encode only uppercase letters, numbers and  [Space], $, %, *, +, -, ., /, :
	AlphaNumeric
)

func (e Encoding) getEncoder() encodeFn {
	switch e {
	case Auto:
		return encodeAuto
	case Numeric:
		return encodeNumeric
	case AlphaNumeric:
		return encodeAlphaNumeric
	}
	return nil
}

func (e Encoding) String() string {
	switch e {
	case Auto:
		return "Auto"
	case Numeric:
		return "Numeric"
	case AlphaNumeric:
		return "AlphaNumeric"
	}
	return ""
}

// Encodes the given content to a QR barcode
func Encode(content string, level ErrorCorrectionLevel, mode Encoding) (barcode.Barcode, error) {
	bits, vi, err := mode.getEncoder()(content, level)
	if err != nil {
		return nil, err
	}
	if bits == nil || vi == nil {
		return nil, fmt.Errorf("Unable to encode \"%s\" with error correction level %s and encoding mode %s", content, level, mode)
	}

	blocks := splitToBlocks(bits.ItterateBytes(), vi)
	data := blocks.interleave(vi)
	result := render(data, vi)
	result.content = content
	return result, nil
}

func render(data []byte, vi *versionInfo) *qrcode {
	dim := vi.modulWidth()
	results := make([]*qrcode, 8)
	for i := 0; i < 8; i++ {
		results[i] = newBarcode(dim)
	}

	occupied := newBarcode(dim)

	setAll := func(x int, y int, val bool) {
		occupied.Set(x, y, true)
		for i := 0; i < 8; i++ {
			results[i].Set(x, y, val)
		}
	}

	drawFinderPatterns(vi, setAll)
	drawAlignmentPatterns(occupied, vi, setAll)

	//Timing Pattern:
	var i int
	for i = 0; i < dim; i++ {
		if !occupied.Get(i, 6) {
			setAll(i, 6, i%2 == 0)
		}
		if !occupied.Get(6, i) {
			setAll(6, i, i%2 == 0)
		}
	}
	// Dark Module
	setAll(8, dim-8, true)

	drawVersionInfo(vi, setAll)
	drawFormatInfo(vi, -1, occupied.Set)
	for i := 0; i < 8; i++ {
		drawFormatInfo(vi, i, results[i].Set)
	}

	// Write the data
	var curBitNo int = 0

	for pos := range itterateModules(occupied) {
		if curBitNo < len(data)*8 {
			curBit := ((data[curBitNo/8] >> uint(7-(curBitNo%8))) & 1) == 1
			for i := 0; i < 8; i++ {
				setMasked(pos.X, pos.Y, curBit, i, results[i].Set)
			}
			curBitNo += 1
		}
	}

	lowestPenalty := ^uint(0)
	lowestPenaltyIdx := -1
	for i := 0; i < 8; i++ {
		p := results[i].calcPenalty()
		if p < lowestPenalty {
			lowestPenalty = p
			lowestPenaltyIdx = i
		}
	}
	return results[lowestPenaltyIdx]
}

func setMasked(x, y int, val bool, mask int, set func(int, int, bool)) {
	switch mask {
	case 0:
		val = val != (((y + x) % 2) == 0)
		break
	case 1:
		val = val != ((y % 2) == 0)
		break
	case 2:
		val = val != ((x % 3) == 0)
		break
	case 3:
		val = val != (((y + x) % 3) == 0)
		break
	case 4:
		val = val != (((y/2 + x/3) % 2) == 0)
		break
	case 5:
		val = val != (((y*x)%2)+((y*x)%3) == 0)
		break
	case 6:
		val = val != ((((y*x)%2)+((y*x)%3))%2 == 0)
		break
	case 7:
		val = val != ((((y+x)%2)+((y*x)%3))%2 == 0)
	}
	set(x, y, val)
}

func itterateModules(occupied *qrcode) <-chan image.Point {
	result := make(chan image.Point)
	allPoints := make(chan image.Point)
	go func() {
		curX := occupied.dimension - 1
		curY := occupied.dimension - 1
		isUpward := true

		for true {
			if isUpward {
				allPoints <- image.Pt(curX, curY)
				allPoints <- image.Pt(curX-1, curY)
				curY -= 1
				if curY < 0 {
					curY = 0
					curX -= 2
					if curX == 6 {
						curX -= 1
					}
					if curX < 0 {
						break
					}
					isUpward = false
				}
			} else {
				allPoints <- image.Pt(curX, curY)
				allPoints <- image.Pt(curX-1, curY)
				curY += 1
				if curY >= occupied.dimension {
					curY = occupied.dimension - 1
					curX -= 2
					if curX == 6 {
						curX -= 1
					}
					isUpward = true
					if curX < 0 {
						break
					}
				}
			}
		}

		close(allPoints)
	}()
	go func() {
		for pt := range allPoints {
			if !occupied.Get(pt.X, pt.Y) {
				result <- pt
			}
		}
		close(result)
	}()
	return result
}

func drawFinderPatterns(vi *versionInfo, set func(int, int, bool)) {
	dim := vi.modulWidth()
	drawPattern := func(xoff int, yoff int) {
		for x := -1; x < 8; x++ {
			for y := -1; y < 8; y++ {
				val := (x == 0 || x == 6 || y == 0 || y == 6 || (x > 1 && x < 5 && y > 1 && y < 5)) && (x <= 6 && y <= 6 && x >= 0 && y >= 0)

				if x+xoff >= 0 && x+xoff < dim && y+yoff >= 0 && y+yoff < dim {
					set(x+xoff, y+yoff, val)
				}
			}
		}
	}
	drawPattern(0, 0)
	drawPattern(0, dim-7)
	drawPattern(dim-7, 0)
}

func drawAlignmentPatterns(occupied *qrcode, vi *versionInfo, set func(int, int, bool)) {
	drawPattern := func(xoff int, yoff int) {
		for x := -2; x <= 2; x++ {
			for y := -2; y <= 2; y++ {
				val := x == -2 || x == 2 || y == -2 || y == 2 || (x == 0 && y == 0)
				set(x+xoff, y+yoff, val)
			}
		}
	}
	positions := vi.alignmentPatternPlacements()

	for _, x := range positions {
		for _, y := range positions {
			if occupied.Get(x, y) {
				continue
			}
			drawPattern(x, y)
		}
	}
}

func drawFormatInfo(vi *versionInfo, usedMask int, set func(int, int, bool)) {
	var formatInfo []bool
	switch vi.Level {
	case L:
		switch usedMask {
		case 0:
			formatInfo = []bool{true, true, true, false, true, true, true, true, true, false, false, false, true, false, false}
			break
		case 1:
			formatInfo = []bool{true, true, true, false, false, true, false, true, true, true, true, false, false, true, true}
			break
		case 2:
			formatInfo = []bool{true, true, true, true, true, false, true, true, false, true, false, true, false, true, false}
			break
		case 3:
			formatInfo = []bool{true, true, true, true, false, false, false, true, false, false, true, true, true, false, true}
			break
		case 4:
			formatInfo = []bool{true, true, false, false, true, true, false, false, false, true, false, true, true, true, true}
			break
		case 5:
			formatInfo = []bool{true, true, false, false, false, true, true, false, false, false, true, true, false, false, false}
			break
		case 6:
			formatInfo = []bool{true, true, false, true, true, false, false, false, true, false, false, false, false, false, true}
			break
		case 7:
			formatInfo = []bool{true, true, false, true, false, false, true, false, true, true, true, false, true, true, false}
			break
		}
		break
	case M:
		switch usedMask {
		case 0:
			formatInfo = []bool{true, false, true, false, true, false, false, false, false, false, true, false, false, true, false}
			break
		case 1:
			formatInfo = []bool{true, false, true, false, false, false, true, false, false, true, false, false, true, false, true}
			break
		case 2:
			formatInfo = []bool{true, false, true, true, true, true, false, false, true, true, true, true, true, false, false}
			break
		case 3:
			formatInfo = []bool{true, false, true, true, false, true, true, false, true, false, false, true, false, true, true}
			break
		case 4:
			formatInfo = []bool{true, false, false, false, true, false, true, true, true, true, true, true, false, false, true}
			break
		case 5:
			formatInfo = []bool{true, false, false, false, false, false, false, true, true, false, false, true, true, true, false}
			break
		case 6:
			formatInfo = []bool{true, false, false, true, true, true, true, true, false, false, true, false, true, true, true}
			break
		case 7:
			formatInfo = []bool{true, false, false, true, false, true, false, true, false, true, false, false, false, false, false}
			break
		}
		break
	case Q:
		switch usedMask {
		case 0:
			formatInfo = []bool{false, true, true, false, true, false, true, false, true, false, true, true, true, true, true}
			break
		case 1:
			formatInfo = []bool{false, true, true, false, false, false, false, false, true, true, false, true, false, false, false}
			break
		case 2:
			formatInfo = []bool{false, true, true, true, true, true, true, false, false, true, true, false, false, false, true}
			break
		case 3:
			formatInfo = []bool{false, true, true, true, false, true, false, false, false, false, false, false, true, true, false}
			break
		case 4:
			formatInfo = []bool{false, true, false, false, true, false, false, true, false, true, true, false, true, false, false}
			break
		case 5:
			formatInfo = []bool{false, true, false, false, false, false, true, true, false, false, false, false, false, true, true}
			break
		case 6:
			formatInfo = []bool{false, true, false, true, true, true, false, true, true, false, true, true, false, true, false}
			break
		case 7:
			formatInfo = []bool{false, true, false, true, false, true, true, true, true, true, false, true, true, false, true}
			break
		}
		break
	case H:
		switch usedMask {
		case 0:
			formatInfo = []bool{false, false, true, false, true, true, false, true, false, false, false, true, false, false, true}
			break
		case 1:
			formatInfo = []bool{false, false, true, false, false, true, true, true, false, true, true, true, true, true, false}
			break
		case 2:
			formatInfo = []bool{false, false, true, true, true, false, false, true, true, true, false, false, true, true, true}
			break
		case 3:
			formatInfo = []bool{false, false, true, true, false, false, true, true, true, false, true, false, false, false, false}
			break
		case 4:
			formatInfo = []bool{false, false, false, false, true, true, true, false, true, true, false, false, false, true, false}
			break
		case 5:
			formatInfo = []bool{false, false, false, false, false, true, false, false, true, false, true, false, true, false, true}
			break
		case 6:
			formatInfo = []bool{false, false, false, true, true, false, true, false, false, false, false, true, true, false, false}
			break
		case 7:
			formatInfo = []bool{false, false, false, true, false, false, false, false, false, true, true, true, false, true, true}
			break
		}
		break
	}

	if usedMask == -1 {
		formatInfo = []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true} // Set all to true cause -1 --> occupied mask.
	}

	if len(formatInfo) == 15 {
		dim := vi.modulWidth()
		set(0, 8, formatInfo[0])
		set(1, 8, formatInfo[1])
		set(2, 8, formatInfo[2])
		set(3, 8, formatInfo[3])
		set(4, 8, formatInfo[4])
		set(5, 8, formatInfo[5])
		set(7, 8, formatInfo[6])
		set(8, 8, formatInfo[7])
		set(8, 7, formatInfo[8])
		set(8, 5, formatInfo[9])
		set(8, 4, formatInfo[10])
		set(8, 3, formatInfo[11])
		set(8, 2, formatInfo[12])
		set(8, 1, formatInfo[13])
		set(8, 0, formatInfo[14])

		set(8, dim-1, formatInfo[0])
		set(8, dim-2, formatInfo[1])
		set(8, dim-3, formatInfo[2])
		set(8, dim-4, formatInfo[3])
		set(8, dim-5, formatInfo[4])
		set(8, dim-6, formatInfo[5])
		set(8, dim-7, formatInfo[6])
		set(dim-8, 8, formatInfo[7])
		set(dim-7, 8, formatInfo[8])
		set(dim-6, 8, formatInfo[9])
		set(dim-5, 8, formatInfo[10])
		set(dim-4, 8, formatInfo[11])
		set(dim-3, 8, formatInfo[12])
		set(dim-2, 8, formatInfo[13])
		set(dim-1, 8, formatInfo[14])
	}
}

func drawVersionInfo(vi *versionInfo, set func(int, int, bool)) {
	var versionInfoBits []bool

	switch vi.Version {
	case 7:
		versionInfoBits = []bool{false, false, false, true, true, true, true, true, false, false, true, false, false, true, false, true, false, false}
		break
	case 8:
		versionInfoBits = []bool{false, false, true, false, false, false, false, true, false, true, true, false, true, true, true, true, false, false}
		break
	case 9:
		versionInfoBits = []bool{false, false, true, false, false, true, true, false, true, false, true, false, false, true, true, false, false, true}
		break
	case 10:
		versionInfoBits = []bool{false, false, true, false, true, false, false, true, false, false, true, true, false, true, false, false, true, true}
		break
	case 11:
		versionInfoBits = []bool{false, false, true, false, true, true, true, false, true, true, true, true, true, true, false, true, true, false}
		break
	case 12:
		versionInfoBits = []bool{false, false, true, true, false, false, false, true, true, true, false, true, true, false, false, false, true, false}
		break
	case 13:
		versionInfoBits = []bool{false, false, true, true, false, true, true, false, false, false, false, true, false, false, false, true, true, true}
		break
	case 14:
		versionInfoBits = []bool{false, false, true, true, true, false, false, true, true, false, false, false, false, false, true, true, false, true}
		break
	case 15:
		versionInfoBits = []bool{false, false, true, true, true, true, true, false, false, true, false, false, true, false, true, false, false, false}
		break
	case 16:
		versionInfoBits = []bool{false, true, false, false, false, false, true, false, true, true, false, true, true, true, true, false, false, false}
		break
	case 17:
		versionInfoBits = []bool{false, true, false, false, false, true, false, true, false, false, false, true, false, true, true, true, false, true}
		break
	case 18:
		versionInfoBits = []bool{false, true, false, false, true, false, true, false, true, false, false, false, false, true, false, true, true, true}
		break
	case 19:
		versionInfoBits = []bool{false, true, false, false, true, true, false, true, false, true, false, false, true, true, false, false, true, false}
		break
	case 20:
		versionInfoBits = []bool{false, true, false, true, false, false, true, false, false, true, true, false, true, false, false, true, true, false}
		break
	case 21:
		versionInfoBits = []bool{false, true, false, true, false, true, false, true, true, false, true, false, false, false, false, false, true, true}
		break
	case 22:
		versionInfoBits = []bool{false, true, false, true, true, false, true, false, false, false, true, true, false, false, true, false, false, true}
		break
	case 23:
		versionInfoBits = []bool{false, true, false, true, true, true, false, true, true, true, true, true, true, false, true, true, false, false}
		break
	case 24:
		versionInfoBits = []bool{false, true, true, false, false, false, true, true, true, false, true, true, false, false, false, true, false, false}
		break
	case 25:
		versionInfoBits = []bool{false, true, true, false, false, true, false, false, false, true, true, true, true, false, false, false, false, true}
		break
	case 26:
		versionInfoBits = []bool{false, true, true, false, true, false, true, true, true, true, true, false, true, false, true, false, true, true}
		break
	case 27:
		versionInfoBits = []bool{false, true, true, false, true, true, false, false, false, false, true, false, false, false, true, true, true, false}
		break
	case 28:
		versionInfoBits = []bool{false, true, true, true, false, false, true, true, false, false, false, false, false, true, true, false, true, false}
		break
	case 29:
		versionInfoBits = []bool{false, true, true, true, false, true, false, false, true, true, false, false, true, true, true, true, true, true}
		break
	case 30:
		versionInfoBits = []bool{false, true, true, true, true, false, true, true, false, true, false, true, true, true, false, true, false, true}
		break
	case 31:
		versionInfoBits = []bool{false, true, true, true, true, true, false, false, true, false, false, true, false, true, false, false, false, false}
		break
	case 32:
		versionInfoBits = []bool{true, false, false, false, false, false, true, false, false, true, true, true, false, true, false, true, false, true}
		break
	case 33:
		versionInfoBits = []bool{true, false, false, false, false, true, false, true, true, false, true, true, true, true, false, false, false, false}
		break
	case 34:
		versionInfoBits = []bool{true, false, false, false, true, false, true, false, false, false, true, false, true, true, true, false, true, false}
		break
	case 35:
		versionInfoBits = []bool{true, false, false, false, true, true, false, true, true, true, true, false, false, true, true, true, true, true}
		break
	case 36:
		versionInfoBits = []bool{true, false, false, true, false, false, true, false, true, true, false, false, false, false, true, false, true, true}
		break
	case 37:
		versionInfoBits = []bool{true, false, false, true, false, true, false, true, false, false, false, false, true, false, true, true, true, false}
		break
	case 38:
		versionInfoBits = []bool{true, false, false, true, true, false, true, false, true, false, false, true, true, false, false, true, false, false}
		break
	case 39:
		versionInfoBits = []bool{true, false, false, true, true, true, false, true, false, true, false, true, false, false, false, false, false, true}
		break
	case 40:
		versionInfoBits = []bool{true, false, true, false, false, false, true, true, false, false, false, true, true, false, true, false, false, true}
		break
	}

	if len(versionInfoBits) > 0 {
		for i := 0; i < len(versionInfoBits); i++ {
			x := (vi.modulWidth() - 11) + i%3
			y := i / 3
			set(x, y, versionInfoBits[len(versionInfoBits)-i-1])
			set(y, x, versionInfoBits[len(versionInfoBits)-i-1])
		}
	}

}

func addPaddingAndTerminator(bl *barcode.BitList, vi *versionInfo) {
	for i := 0; i < 4 && bl.Len() < vi.totalDataBytes()*8; i++ {
		bl.AddBit(false)
	}

	for bl.Len()%8 != 0 {
		bl.AddBit(false)
	}

	for i := 0; bl.Len() < vi.totalDataBytes()*8; i++ {
		if i%2 == 0 {
			bl.AddByte(236)
		} else {
			bl.AddByte(17)
		}
	}
}
