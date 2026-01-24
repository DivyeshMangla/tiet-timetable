package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type BlockClassReader struct{}

func (r *BlockClassReader) Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo) {
	subject := getValidCell(file, sheetName, row, col)
	if subject == "" {
		return false, nil
	}

	parsed, ok := parseSubject(subject)
	if !ok {
		return false, nil
	}

	room := getValidCell(file, sheetName, row+1, col)
	if room == "" {
		return false, nil
	}

	// Check for continuation rows to identify block class
	hasContinuation := getValidCell(file, sheetName, row+2, col) != "" ||
		getValidCell(file, sheetName, row+3, col) != ""
	if !hasContinuation {
		return false, nil
	}

	// Prefer teacher from row+3, fallback to row+2
	teacher := getValidCell(file, sheetName, row+3, col)
	if teacher == "" {
		teacher = getValidCell(file, sheetName, row+2, col)
	}
	if teacher == "" {
		return false, nil
	}

	return true, &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(cleanCell(room)),
		Teacher:     types.Teacher(cleanCell(teacher)),
		IsBlock:     true,
	}
}
