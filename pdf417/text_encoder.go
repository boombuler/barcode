package pdf417

import "fmt"

// Since each code word consists of 2 characters, a padding value is
// needed when encoding a single character. 29 is used as padding because
// it's a switch in all 4 submodes, and doesn't add any data.
const padding_value = 29

type subMode byte

const (
	subUpper subMode = iota
	subLower
	subMixed
	subPunct
)

const (
	switchUpper       = '\u00f1'
	switchUpperSingle = '\u00f2'
	switchLower       = '\u00f3'
	switchMixed       = '\u00f4'
	switchPunct       = '\u00f5'
	switchPunctSingle = '\u00f6'
)

type textEncoder struct {
	Switching     map[subMode]map[subMode][]rune
	SwitchSubmode map[rune]subMode
	ReverseLookup map[rune]map[subMode]int
}

func newTextEncoder() *textEncoder {
	encoder := new(textEncoder)

	characterTables := map[subMode][]rune{
		subUpper: []rune{
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
			'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
			'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', ' ',
			switchLower,
			switchMixed,
			switchPunctSingle,
		},

		subLower: []rune{
			'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
			'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
			's', 't', 'u', 'v', 'w', 'x', 'y', 'z', ' ',
			switchUpperSingle,
			switchMixed,
			switchPunctSingle,
		},

		subMixed: []rune{
			'0', '1', '2', '3', '4', '5', '6', '7', '8',
			'9', '&', '\r', '\t', ',', ':', '#', '-', '.',
			'$', '/', '+', '%', '*', '=', '^',
			switchPunct, ' ',
			switchLower,
			switchUpper,
			switchPunctSingle,
		},

		subPunct: []rune{
			';', '<', '>', '@', '[', '\\', ']', '_', '`',
			'~', '!', '\r', '\t', ',', ':', '\n', '-', '.',
			'$', '/', 'g', '|', '*', '(', ')', '?', '{', '}', '\'',
			switchUpper,
		},
	}

	encoder.Switching = map[subMode]map[subMode][]rune{
		subUpper: map[subMode][]rune{
			subLower: []rune{switchLower},
			subMixed: []rune{switchMixed},
			subPunct: []rune{switchMixed, switchPunct},
		},

		subLower: map[subMode][]rune{
			subUpper: []rune{switchMixed, switchUpper},
			subMixed: []rune{switchMixed},
			subPunct: []rune{switchMixed, switchPunct},
		},

		subMixed: map[subMode][]rune{
			subUpper: []rune{switchUpper},
			subLower: []rune{switchLower},
			subPunct: []rune{switchPunct},
		},

		subPunct: map[subMode][]rune{
			subUpper: []rune{switchUpper},
			subLower: []rune{switchUpper, switchLower},
			subMixed: []rune{switchUpper, switchMixed},
		},
	}

	encoder.SwitchSubmode = map[rune]subMode{
		switchUpper: subUpper,
		switchLower: subLower,
		switchPunct: subPunct,
		switchMixed: subMixed,
	}

	encoder.ReverseLookup = make(map[rune]map[subMode]int)
	for submode, codes := range characterTables {
		for row, char := range codes {
			if encoder.ReverseLookup[char] == nil {
				encoder.ReverseLookup[char] = make(map[subMode]int)
			}

			encoder.ReverseLookup[char][submode] = int(row)
		}
	}

	return encoder
}

func (encoder *textEncoder) CanEncode(char rune) bool {
	switch char {
	case switchUpper,
		switchUpperSingle,
		switchLower,
		switchMixed,
		switchPunct,
		switchPunctSingle:
		return false
	default:
		return encoder.ReverseLookup[char] != nil
	}
}

func (*textEncoder) GetSwitchCode(data string) int {
	return switchCodeText
}

func (encoder *textEncoder) Encode(data string) ([]int, error) {
	interim, err := encoder.encodeInterim(data)
	if err != nil {
		return interim, err
	}

	return encoder.encodeFinal(interim), nil
}

func (encoder *textEncoder) encodeInterim(data string) ([]int, error) {
	submode := subUpper

	codes := []int{}
	var err error
	for _, char := range data {
		if !encoder.existsInSubmode(char, submode) {
			prevSubmode := submode

			submode, err = encoder.getSubmode(char)
			if err != nil {
				return codes, err
			}

			switchCodes := encoder.getSwitchCodes(prevSubmode, submode)

			codes = append(codes, switchCodes...)
		}

		codes = append(
			codes,
			encoder.getCharacterCode(char, submode),
		)
	}

	return codes, nil
}

func (encoder *textEncoder) getSubmode(char rune) (subMode, error) {
	if lookup, ok := encoder.ReverseLookup[char]; ok {
		for key := range lookup {
			return key, nil
		}
	}
	return subLower, fmt.Errorf("unable to find submode for %q", char)
}

func (encoder *textEncoder) getSwitchCodes(from, to subMode) []int {
	switches := encoder.Switching[from][to]

	codes := []int{}

	for _, switcher := range switches {
		codes = append(codes, encoder.getCharacterCode(switcher, from))

		from = encoder.SwitchSubmode[switcher]
	}

	return codes
}

func (*textEncoder) encodeFinal(codes []int) []int {
	codeWords := []int{}

	chunks := [][]int{}
	chunkPart := []int{}
	i := 0
	for _, code := range codes {
		chunkPart = append(chunkPart, code)

		i++

		if i%2 == 0 {
			chunks = append(chunks, chunkPart)

			chunkPart = []int{}
		}
	}

	if len(chunkPart) > 0 {
		chunks = append(chunks, chunkPart)
	}

	for _, chunk := range chunks {
		if len(chunk) == 1 {
			chunk = append(chunk, padding_value)
		}

		codeWords = append(
			codeWords,
			30*chunk[0]+chunk[1],
		)
	}

	return codeWords
}

func (encoder *textEncoder) getCharacterCode(char rune, submode subMode) int {
	cw, ok := encoder.ReverseLookup[char][submode]

	if !ok {
		panic("This is not possible")
	}

	return cw
}

func (encoder *textEncoder) existsInSubmode(char rune, submode subMode) bool {
	_, ok := encoder.ReverseLookup[char][submode]

	return ok
}
