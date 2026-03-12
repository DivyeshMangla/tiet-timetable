package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type LargeBlockClassReader struct{}

func (r LargeBlockClassReader) Read(ws *excel.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool) {
	region, found := ws.HorizontalMergedRegion(row, col)
	if !found || !isWideEnough(region) {
		return nil, false
	}

	startRow := region.StartRow
	startCol := region.StartCol

	subjectValue, err := ws.Cell(startRow, startCol)
	if err != nil {
		return nil, false
	}

	subjectMatcher := utils.NewValueMatcher(CleanCell(subjectValue), subjectCodePattern)
	if !subjectMatcher.Valid() || !subjectMatcher.HasOneValue() {
		return nil, false
	}

	roomValue, err := ws.Cell(startRow+1, startCol)
	if err != nil {
		return nil, false
	}

	room := CleanCell(roomValue)
	if room == "" {
		return nil, false
	}

	cont1, _ := ws.Cell(startRow+2, startCol)
	cont2, _ := ws.Cell(startRow+3, startCol)

	cont1Clean := CleanCell(cont1)
	cont2Clean := CleanCell(cont2)

	if cont1Clean == "" && cont2Clean == "" {
		return nil, false
	}

	// If row+2 looks like a subject code, it's the next time slot, not a continuation
	if cont1Clean != "" {
		contMatcher := utils.NewValueMatcher(cont1Clean, subjectCodePattern)
		if contMatcher.Valid() {
			return nil, false
		}
	}

	teacherValue := CleanCell(cont2)
	if teacherValue == "" {
		teacherValue = CleanCell(cont1)
	}

	if teacherValue == "" {
		return nil, false
	}

	code, ct := parseSubjectCode(subjectMatcher.Values()[0])

	class := types.Class{
		SubjectCode: code,
		ClassType:   ct,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacherValue),
	}

	return &types.ClassSlot{
		Start: start,
		End:   start + 1,
		Classes: []types.Class{
			class,
		},
	}, true
}
