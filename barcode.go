package barcode

import "image"

type Metadata struct {
	CodeKind   string
	Dimensions byte
}

type Barcode interface {
	image.Image
	Metadata() Metadata
	Content() string
}
