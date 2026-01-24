package image

import (
	"fmt"
	"image/color"
)

func HexToRGBA(hex string) color.RGBA {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
