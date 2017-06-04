package pdf417

import (
	"fmt"
)

const (
	switchCodeText   = 900
	switchCodeNumber = 902
)

type encoder interface {
	CanEncode(char rune) bool
	GetSwitchCode(data string) int
	Encode(data string) ([]int, error)
}

type dataEncoder struct {
	Encoders       []encoder
	DefaultEncoder encoder
}

type chain struct {
	Data    string
	Encoder encoder
}

func (c chain) Encode() ([]int, error) {
	return c.Encoder.Encode(c.Data)
}

func (c chain) GetSwitchCode() int {
	return c.Encoder.GetSwitchCode(c.Data)
}

func newDataEncoder() *dataEncoder {
	textEncoder := newTextEncoder()

	encoder := &dataEncoder{
		[]encoder{
			newNumberEncoder(),
			textEncoder,
		},
		textEncoder,
	}

	return encoder
}

func (dataEncoder *dataEncoder) Encode(data string) ([]int, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("Nothing to encode")
	}

	chains, err := dataEncoder.SplitToChains(data)
	if err != nil {
		return nil, err
	}

	if len(chains) == 0 {
		return nil, fmt.Errorf("%q can not be encoded!", data)
	}

	currentSwitchCode := switchCodeText

	cws := []int{}

	for _, chain := range chains {
		encoded, err := chain.Encode()
		if err != nil {
			return cws, err
		}

		if switchCode := chain.GetSwitchCode(); currentSwitchCode != switchCode {
			cws = append(cws, switchCode)
			currentSwitchCode = switchCode
		}

		cws = append(cws, encoded...)
	}

	return cws, nil
}

func (dataEncoder *dataEncoder) SplitToChains(data string) ([]chain, error) {
	chains := []chain{}
	chainData := ""
	encoder := dataEncoder.DefaultEncoder

	for _, char := range data {
		newEncoder, err := dataEncoder.getEncoder(char)
		if err != nil {
			return nil, err
		}

		if newEncoder != encoder {
			if len(chainData) > 0 {
				chains = append(chains, chain{chainData, encoder})
				chainData = ""
			}

			encoder = newEncoder
		}

		chainData = chainData + string(char)
	}

	if len(chainData) > 0 {
		chains = append(chains, chain{chainData, encoder})
	}

	return chains, nil
}

func (dataEncoder *dataEncoder) getEncoder(char rune) (encoder, error) {
	for _, encoder := range dataEncoder.Encoders {
		if encoder.CanEncode(char) {
			return encoder, nil
		}
	}
	return nil, fmt.Errorf("Cannot encode character %q", string(char))
}
