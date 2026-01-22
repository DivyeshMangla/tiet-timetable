package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type SingleClassReader struct{}

func (r *SingleClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	subjectValue, err := parser.GetCell(file, sheetName, row, col)
	if err != nil {
		return false
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(subjectValue))
	if !ok || !parser.IsSubjectCode(string(parsed.Code)) {
		return false
	}

	roomValue, err := parser.GetCell(file, sheetName, row+1, col)
	if err != nil {
		return false
	}

	teacherValue, err := parser.GetCell(file, sheetName, row+1, col+1)
	if err != nil {
		return false
	}

	if !isValid(roomValue) || !isValid(teacherValue) {
		return false
	}

	return true
}

func (r *SingleClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subjectValue, err := parser.GetCell(file, sheetName, row, col)
	if err != nil {
		return nil
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(subjectValue))
	if !ok {
		return nil
	}

	roomValue, _ := parser.GetCell(file, sheetName, row+1, col)
	teacherValue, _ := parser.GetCell(file, sheetName, row+1, col+1)

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(parser.GetCellString(roomValue)),
		Teacher:     types.Teacher(parser.GetCellString(teacherValue)),
		IsBlock:     false,
	}
}
