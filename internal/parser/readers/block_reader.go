package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type BlockClassReader struct{}

func (r BlockClassReader) Read(ws *parser.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool) {
	subjectValue, err := ws.Cell(row, col)
	if err != nil {
		return nil, false
	}

	subjectMatcher := utils.NewValueMatcher(CleanCell(subjectValue), subjectCodePattern)
	if !subjectMatcher.Valid() || !subjectMatcher.HasOneValue() {
		return nil, false
	}

	roomValue, err := ws.Cell(row+1, col)
	if err != nil {
		return nil, false
	}

	roomMatcher := utils.NewValueMatcher(CleanCell(roomValue), nil)
	if !roomMatcher.HasOneValue() {
		return nil, false
	}

	// Detect continuation rows (block indicator)
	cont1, _ := ws.Cell(row+2, col)
	cont2, _ := ws.Cell(row+3, col)

	if CleanCell(cont1) == "" && CleanCell(cont2) == "" {
		return nil, false
	}

	// Prefer teacher from row+3, fallback to row+2
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
		End:   start + 1, // block spans two slots
		Classes: []types.Class{
			class,
		},
	}, true
}
