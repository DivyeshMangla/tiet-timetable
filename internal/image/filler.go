package image

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	BackgroundColorR = 22
	BackgroundColorG = 24
	BackgroundColorB = 25
	ColorTolerance   = 10
)

var BackgroundColor = color.RGBA{R: BackgroundColorR, G: BackgroundColorG, B: BackgroundColorB, A: 255}

type CapsuleFiller struct {
	img *image.RGBA
}

func NewCapsuleFiller(imgPath string) (*CapsuleFiller, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	return &CapsuleFiller{img: rgba}, nil
}

func (cf *CapsuleFiller) FillCell(timeSlot model.TimeSlot, day model.Day, fillColor color.RGBA) {
	cell := GetCell(timeSlot, day)

	for y := cell.Y; y < cell.Y+cell.Height; y++ {
		for x := cell.X; x < cell.X+cell.Width; x++ {
			if x >= cf.img.Bounds().Max.X || y >= cf.img.Bounds().Max.Y {
				continue
			}

			r, g, b, a := cf.img.At(x, y).RGBA()
			pixelColor := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}

			if !isBackgroundColor(pixelColor) {
				cf.img.Set(x, y, fillColor)
			}
		}
	}
}

func isBackgroundColor(c color.RGBA) bool {
	return abs(c.R, BackgroundColorR) <= ColorTolerance &&
		abs(c.G, BackgroundColorG) <= ColorTolerance &&
		abs(c.B, BackgroundColorB) <= ColorTolerance
}

func abs(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

func (cf *CapsuleFiller) Save(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, cf.img)
}
