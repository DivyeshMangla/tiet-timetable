package image

import (
	"image/color"
	"io"

	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type TimetableDrawer struct {
	filler *CapsuleFiller
}

func NewTimetableDrawer() (*TimetableDrawer, error) {
	filler, err := NewCapsuleFiller()
	if err != nil {
		return nil, err
	}

	return &TimetableDrawer{
		filler: filler,
	}, nil
}

func (td *TimetableDrawer) getFillColor(classType types.ClassType) color.RGBA {
	switch classType {
	case types.LECTURE:
		return LectureColor
	case types.TUTORIAL:
		return TutorialColor
	case types.PRACTICAL:
		return PracticalColor
	default:
		return CellColor
	}
}

func (td *TimetableDrawer) drawCells(timetable *types.RenderableTimetable) error {
	for day, infos := range timetable.Days {
		for _, renderInfo := range infos {
			fillColor := td.getFillColor(renderInfo.ClassType)

			if renderInfo.IsBlock() {
				err := td.filler.FillVerticalWithText(
					renderInfo.Start,
					day,
					fillColor,
					renderInfo.Text,
				)
				if err != nil {
					return err
				}
			} else {
				err := td.filler.FillCellWithText(
					renderInfo.Start,
					day,
					fillColor,
					renderInfo.Text,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (td *TimetableDrawer) DrawTimetable(timetable *types.RenderableTimetable, outputPath string) error {
	if err := td.drawCells(timetable); err != nil {
		return err
	}
	return td.filler.Save(outputPath)
}

func (td *TimetableDrawer) WriteTimetable(timetable *types.RenderableTimetable, w io.Writer) error {
	if err := td.drawCells(timetable); err != nil {
		return err
	}
	return td.filler.WriteTo(w)
}
