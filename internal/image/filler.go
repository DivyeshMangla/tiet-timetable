package image

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	ColorTolerance  = 10
	DefaultFontSize = 53.0
)

type CapsuleFiller struct {
	img  *image.RGBA
	font *truetype.Font
}

func NewCapsuleFiller() (*CapsuleFiller, error) {
	bgFile, err := GetBackground(Background)
	if err != nil {
		return nil, err
	}
	defer bgFile.Close()

	img, err := png.Decode(bgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	fontBytes, err := GetFont(FontFile)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return &CapsuleFiller{img: rgba, font: f}, nil
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

			// Skip background pixels
			if isBackgroundColor(pixelColor) {
				continue
			}

			// Calculate darkness relative to cell background
			darknessR := int(CellColor.R) - int(pixelColor.R)
			darknessG := int(CellColor.G) - int(pixelColor.G)
			darknessB := int(CellColor.B) - int(pixelColor.B)

			// Apply the same darkness to the fill color
			newColor := color.RGBA{
				R: clamp(int(fillColor.R) - darknessR),
				G: clamp(int(fillColor.G) - darknessG),
				B: clamp(int(fillColor.B) - darknessB),
				A: fillColor.A,
			}

			cf.img.Set(x, y, newColor)
		}
	}
}

func (cf *CapsuleFiller) FillCellWithText(timeSlot model.TimeSlot, day model.Day, fillColor color.RGBA, text string) error {
	cf.FillCell(timeSlot, day, fillColor)

	cell := GetCell(timeSlot, day)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(cf.font)
	c.SetFontSize(DefaultFontSize)
	c.SetClip(cf.img.Bounds())
	c.SetDst(cf.img)
	c.SetSrc(image.NewUniform(TextColor))
	c.SetHinting(font.HintingFull)

	face := truetype.NewFace(cf.font, &truetype.Options{
		Size:    DefaultFontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	bounds, _ := font.BoundString(face, text)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	leftBearing := bounds.Min.X.Ceil()

	metrics := face.Metrics()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	textHeight := ascent + descent

	x := cell.X + (cell.Width-textWidth)/2 - leftBearing
	y := cell.Y + (cell.Height / 2) - (textHeight / 2) + ascent

	pt := freetype.Pt(x, y)
	_, err := c.DrawString(text, pt)

	return err
}

func isBackgroundColor(c color.RGBA) bool {
	return abs(c.R, BackgroundColor.R) <= ColorTolerance &&
		abs(c.G, BackgroundColor.G) <= ColorTolerance &&
		abs(c.B, BackgroundColor.B) <= ColorTolerance
}

func (cf *CapsuleFiller) Save(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()
	return png.Encode(file, cf.img)
}
