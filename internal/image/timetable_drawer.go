package image

import (
	"fmt"
	"image/color"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
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

func (td *TimetableDrawer) getFillColor(classType model.ClassType) color.RGBA {
	switch classType {
	case model.LECTURE:
		return LectureColor
	case model.TUTORIAL:
		return TutorialColor
	case model.PRACTICAL:
		return PracticalColor
	default:
		return CellColor
	}
}

func (td *TimetableDrawer) DrawTimetable(entries []model.TimetableEntry, outputPath string) error {
	for _, entry := range entries {
		text := fmt.Sprintf("%s - %s", entry.ClassInfo.SubjectCode, entry.ClassInfo.Room)
		fillColor := td.getFillColor(entry.ClassInfo.ClassType)

		if entry.ClassInfo.IsBlock {
			err := td.filler.FillVerticalWithText(
				entry.TimeSlot,
				entry.Day,
				fillColor,
				text,
			)
			if err != nil {
				return err
			}
		} else {
			err := td.filler.FillCellWithText(
				entry.TimeSlot,
				entry.Day,
				fillColor,
				text,
			)
			if err != nil {
				return err
			}
		}
	}

	return td.filler.Save(outputPath)
}
