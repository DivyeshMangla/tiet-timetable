package image

import (
	"image"
	"image/color"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func (cf *CapsuleFiller) FillVertical(ts model.TimeSlot, day model.Day, fillColor color.RGBA) {
	xStart := ScheduleGrid.XBounds[day].Start
	xEnd := ScheduleGrid.XBounds[day].End

	upperYStart := ScheduleGrid.YBounds[ts].Start
	upperYEnd := ScheduleGrid.YBounds[ts].End
	upperBounds := cf.CalculateBounds(upperYStart, upperYEnd, xStart, xEnd)

	lowerYStart := ScheduleGrid.YBounds[ts+1].Start
	lowerYEnd := ScheduleGrid.YBounds[ts+1].End
	lowerBounds := cf.CalculateBounds(lowerYStart, lowerYEnd, xStart, xEnd)

	cf.FillCell(ts, day, fillColor)
	cf.FillCell(ts+1, day, fillColor)

	if upperBounds.Left != -1 && lowerBounds.Left != -1 {
		cf.FillRectangleBetweenBounds(upperBounds, lowerBounds, fillColor)
	}
}

func (cf *CapsuleFiller) FillVerticalWithText(ts model.TimeSlot, day model.Day, fillColor color.RGBA, text string) error {
	cf.FillVertical(ts, day, fillColor)

	cell := GetMergedCell(ts, day)

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

	_, err := c.DrawString(text, freetype.Pt(x, y))
	return err
}
