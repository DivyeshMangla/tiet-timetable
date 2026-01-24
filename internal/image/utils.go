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

func abs(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

func clamp(val int) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}
