package image

import (
	"image/color"
)

type RowBounds struct {
	Left   int
	Right  int
	YLevel int
}

func (cf *CapsuleFiller) CalculateBounds(yStart, yEnd, xStart, xEnd int) RowBounds {
	bounds := RowBounds{Left: -1, Right: -1, YLevel: -1}
	minX := int(^uint(0) >> 1)
	maxX := -1

	for y := yStart; y < yEnd; y++ {
		if y >= cf.img.Bounds().Max.Y {
			continue
		}

		rowMinX := int(^uint(0) >> 1)
		rowMaxX := -1

		for x := xStart; x < xEnd; x++ {
			if x >= cf.img.Bounds().Max.X {
				continue
			}

			r, g, b, a := cf.img.At(x, y).RGBA()
			pixel := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}

			if isBackgroundColor(pixel) {
				continue
			}

			if x < rowMinX {
				rowMinX = x
			}
			if x > rowMaxX {
				rowMaxX = x
			}
		}

		if rowMinX < minX {
			minX = rowMinX
			bounds.Left = minX
			bounds.YLevel = y
		}
		if rowMaxX > maxX {
			maxX = rowMaxX
			bounds.Right = maxX
		}
	}

	return bounds
}

func (cf *CapsuleFiller) FillRectangleBetweenBounds(upperRowBound, lowerRowBound RowBounds, fillColor color.RGBA) {
	left := upperRowBound.Left
	if lowerRowBound.Left < left {
		left = lowerRowBound.Left
	}

	right := upperRowBound.Right
	if lowerRowBound.Right > right {
		right = lowerRowBound.Right
	}

	for y := upperRowBound.YLevel; y <= lowerRowBound.YLevel; y++ {
		if y >= cf.img.Bounds().Max.Y {
			continue
		}

		for x := left; x <= right; x++ {
			if x >= cf.img.Bounds().Max.X {
				continue
			}

			cf.img.Set(x, y, fillColor)
		}
	}
}
