package barcode

// ECI is an Extended Channel Interpretation designator: a number from
// 0 to MaxECI, registered with AIM, that tells a reader how to interpret
// the data following it in a symbol.  Most designators name the character
// set that the bytes are in.
//
// QR Code, Data Matrix, Aztec Code and PDF417 can carry an ECI.
// Their packages each have an EncodeWithECI function that writes one
// at the start of the symbol.  A reader that does not support ECIs
// may show one as data or refuse the symbol.
type ECI int

const (
	// NoECI asks for a symbol that carries no ECI.
	NoECI ECI = -1
	// MaxECI is the largest designator.
	MaxECI ECI = 999999
)

// Designators of the character sets in the AIM ECI register.
const (
	ECICP437       ECI = 2
	ECIISO8859_1   ECI = 3
	ECIISO8859_2   ECI = 4
	ECIISO8859_3   ECI = 5
	ECIISO8859_4   ECI = 6
	ECIISO8859_5   ECI = 7
	ECIISO8859_6   ECI = 8
	ECIISO8859_7   ECI = 9
	ECIISO8859_8   ECI = 10
	ECIISO8859_9   ECI = 11
	ECIISO8859_10  ECI = 12
	ECIISO8859_11  ECI = 13
	ECIISO8859_13  ECI = 15
	ECIISO8859_14  ECI = 16
	ECIISO8859_15  ECI = 17
	ECIISO8859_16  ECI = 18
	ECIShiftJIS    ECI = 20
	ECIWindows1250 ECI = 21
	ECIWindows1251 ECI = 22
	ECIWindows1252 ECI = 23
	ECIWindows1256 ECI = 24
	ECIUTF16BE     ECI = 25
	ECIUTF8        ECI = 26
	ECIASCII       ECI = 27
	ECIBig5        ECI = 28
	ECIGB2312      ECI = 29
	ECIEUCKR       ECI = 30
	ECIGBK         ECI = 31
	ECIGB18030     ECI = 32
	ECIUTF16LE     ECI = 33
	ECIUTF32BE     ECI = 34
	ECIUTF32LE     ECI = 35
)

// Valid reports whether eci can be passed to an EncodeWithECI function.
//
// That is NoECI or a designator from 0 to MaxECI.  PDF417 can carry
// designators up to 811799 only, so it refuses larger ones as well.
func (eci ECI) Valid() bool {
	return eci == NoECI || eci >= 0 && eci <= MaxECI
}
