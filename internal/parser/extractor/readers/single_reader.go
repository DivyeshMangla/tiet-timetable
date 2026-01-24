package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

const (
	singleSubjectRowOffset = 0
	singleRoomRowOffset    = 1
	singleTeacherColOffset = 1
)

type SingleClassReader struct{}

func (r *SingleClassReader) Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo) {
	subject := getValidCell(file, sheetName, row+singleSubjectRowOffset, col)
	if subject == "" {
		return false, nil
	}

	parsed, ok := parseSubject(subject)
	if !ok {
		return false, nil
	}

	room := getValidCell(file, sheetName, row+singleRoomRowOffset, col)
	teacher := getValidCell(file, sheetName, row+singleRoomRowOffset, col+singleTeacherColOffset)
	if room == "" || teacher == "" {
		return false, nil
	}

	regions, err := utils.GetMergedRegions(file, sheetName)
	if err == nil {
		roomRow, roomCol := row+singleRoomRowOffset, col
		teacherRow, teacherCol := row+singleRoomRowOffset, col+singleTeacherColOffset

		for _, region := range regions {
			if region.IsInRange(roomRow, roomCol) && region.StartCol != region.EndCol {
				return false, nil
			}
			if region.IsInRange(teacherRow, teacherCol) && region.StartCol != region.EndCol {
				return false, nil
			}
		}
	}

	return true, &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(cleanCell(room)),
		Teacher:     types.Teacher(cleanCell(teacher)),
		IsBlock:     false,
	}
}
