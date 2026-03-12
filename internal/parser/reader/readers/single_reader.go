package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

const (
	singleSubjectRowOffset = 0
	singleRoomRowOffset    = 1
	singleTeacherColOffset = 1
)

type SingleClassReader struct{}

func (r SingleClassReader) Read(ws *excel.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool) {
	subjectValue, err := ws.Cell(row+singleSubjectRowOffset, col)
	if err != nil {
		return nil, false
	}

	subjectMatcher := utils.NewValueMatcher(CleanCell(subjectValue), subjectCodePattern)
	if !subjectMatcher.Valid() || !subjectMatcher.HasOneValue() {
		return nil, false
	}

	roomValue, err := ws.Cell(row+singleRoomRowOffset, col)
	if err != nil {
		return nil, false
	}

	teacherValue, err := ws.Cell(row+singleRoomRowOffset, col+singleTeacherColOffset)
	if err != nil {
		return nil, false
	}

	roomMatcher := utils.NewValueMatcher(CleanCell(roomValue), nil)
	teacherMatcher := utils.NewValueMatcher(CleanCell(teacherValue), teacherPattern)

	if !teacherMatcher.Valid() || !roomMatcher.HasOneValue() || !teacherMatcher.HasOneValue() {
		return nil, false
	}

	room := roomMatcher.Values()[0]
	teacher := teacherMatcher.Values()[0]
	code, ct := parseSubjectCode(subjectMatcher.Values()[0])

	class := types.Class{
		SubjectCode: code,
		ClassType:   ct,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacher),
	}

	return &types.ClassSlot{
		Start:   start,
		End:     start,
		Classes: []types.Class{class},
	}, true
}
