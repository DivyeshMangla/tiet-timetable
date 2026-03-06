package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type LargeBlockClassReader struct{}

func (r LargeBlockClassReader) Read(ws *parser.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool) {
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

	roomMatcher := utils.NewValueMatcher(CleanCell(roomValue), nil)
	if !roomMatcher.HasOneValue() {
		return nil, false
	}

	cont1, _ := ws.Cell(startRow+2, startCol)
	cont2, _ := ws.Cell(startRow+3, startCol)

	if CleanCell(cont1) == "" && CleanCell(cont2) == "" {
		return nil, false
	}

	teacherValue := CleanCell(cont2)
	if teacherValue == "" {
		teacherValue = CleanCell(cont1)
	}

	teacherMatcher := utils.NewValueMatcher(teacherValue, teacherPattern)
	if !teacherMatcher.Valid() || !teacherMatcher.HasOneValue() {
		return nil, false
	}

	class := types.Class{
		SubjectCode: types.SubjectCode(subjectMatcher.Values()[0]),
		Room:        types.Room(roomMatcher.Values()[0]),
		Teacher:     types.Teacher(teacherMatcher.Values()[0]),
	}

	return &types.ClassSlot{
		Start: start,
		End:   start + 1,
		Classes: []types.Class{
			class,
		},
	}, true
}
