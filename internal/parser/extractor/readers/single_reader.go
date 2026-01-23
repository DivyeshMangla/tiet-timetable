package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

const (
	// singleSubjectRowOffset is the row offset for the subject (0 = same row)
	singleSubjectRowOffset = 0
	// singleRoomRowOffset is the row offset for the room (1 = one row below)
	singleRoomRowOffset = 1
	// singleTeacherColOffset is the column offset for the teacher (1 = one column to the right)
	singleTeacherColOffset = 1
)

// SingleClassReader handles standard single-slot classes where:
//   - Subject is in the current row
//   - Room and teacher are in the row below (room in same column, teacher one column right)
//   - Room and teacher cells must not be horizontally merged
type SingleClassReader struct{}

func (r *SingleClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	subject := getValidCell(file, sheetName, row+singleSubjectRowOffset, col)
	if subject == "" {
		return false
	}

	_, ok := parseSubject(subject)
	if !ok {
		return false
	}

	room := getValidCell(file, sheetName, row+singleRoomRowOffset, col)
	teacher := getValidCell(file, sheetName, row+singleRoomRowOffset, col+singleTeacherColOffset)
	if room == "" || teacher == "" {
		return false
	}

	// Reject if room or teacher cells are horizontally merged
	// (horizontal merges indicate a large class, not a single class)
	regions, err := utils.GetMergedRegions(file, sheetName)
	if err == nil {
		roomRow, roomCol := row+singleRoomRowOffset, col
		teacherRow, teacherCol := row+singleRoomRowOffset, col+singleTeacherColOffset

		for _, region := range regions {
			if region.IsInRange(roomRow, roomCol) && region.StartCol != region.EndCol {
				return false
			}
			if region.IsInRange(teacherRow, teacherCol) && region.StartCol != region.EndCol {
				return false
			}
		}
	}

	return true
}

func (r *SingleClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subject := getValidCell(file, sheetName, row+singleSubjectRowOffset, col)
	if subject == "" {
		return nil
	}

	parsed, ok := parseSubject(subject)
	if !ok {
		return nil
	}

	room := getValidCell(file, sheetName, row+singleRoomRowOffset, col)
	teacher := getValidCell(file, sheetName, row+singleRoomRowOffset, col+singleTeacherColOffset)
	if room == "" || teacher == "" {
		return nil
	}

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(cleanCell(room)),
		Teacher:     types.Teacher(cleanCell(teacher)),
		IsBlock:     false,
	}
}
