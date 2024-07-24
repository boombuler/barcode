package utils

import "image/color"

func ColorModel(depth int) color.Model {
	switch depth {
	case 8:
		return color.GrayModel
	case 24, 32:
		return color.RGBAModel
	default:
		return color.Gray16Model
	}
}

func WhiteColor(depth int) color.Color {
	switch depth {
	case 8:
		return color.Gray{Y: 255}
	case 24, 32:
		return color.RGBA{255, 255, 255, 255}
	default:
		return color.White
	}
}

func BlackColor(depth int) color.Color {
	switch depth {
	case 8:
		return color.Gray{Y: 0}
	case 24, 32:
		return color.RGBA{0, 0, 0, 255}
	default:
		return color.Black
	}
}
