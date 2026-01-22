package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type BlockClassReader struct{}

func (r *BlockClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	// Subject must exist and be a valid subject code
	subject := getValidCell(file, sheetName, row, col)
	if subject == "" {
		return false
	}
	if _, ok := parseSubject(subject); !ok {
		return false
	}

	// Room must exist immediately below
	if getValidCell(file, sheetName, row+1, col) == "" {
		return false
	}

	// At least one continuation row must exist (block indicator)
	return getValidCell(file, sheetName, row+2, col) != "" ||
		getValidCell(file, sheetName, row+3, col) != ""
}

func (r *BlockClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subject := getValidCell(file, sheetName, row, col)
	if subject == "" {
		return nil
	}

	parsed, ok := parseSubject(subject)
	if !ok {
		return nil
	}

	room := getValidCell(file, sheetName, row+1, col)
	if room == "" {
		return nil
	}

	// Prefer teacher from row+3, fallback to row+2
	teacher := getValidCell(file, sheetName, row+3, col)
	if teacher == "" {
		teacher = getValidCell(file, sheetName, row+2, col)
	}
	if teacher == "" {
		return nil
	}

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(cleanCell(room)),
		Teacher:     types.Teacher(cleanCell(teacher)),
		IsBlock:     true,
	}
}
