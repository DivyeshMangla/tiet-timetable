package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type BlockClassReader struct{}

func (r *BlockClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	subjectValue, err := parser.GetCell(file, sheetName, row, col)
	if err != nil || !isValid(subjectValue) {
		return false
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(subjectValue))
	if !ok || !parser.IsSubjectCode(string(parsed.Code)) {
		return false
	}

	roomValue, err := parser.GetCell(file, sheetName, row+1, col)
	if err != nil || !isValid(roomValue) {
		return false
	}

	row2Value, _ := parser.GetCell(file, sheetName, row+2, col)
	row3Value, _ := parser.GetCell(file, sheetName, row+3, col)

	if isValid(row2Value) && isValid(row3Value) {
		return true
	}

	return isValid(row2Value) && !isValid(row3Value)
}

func (r *BlockClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subjectValue, err := parser.GetCell(file, sheetName, row, col)
	if err != nil || !isValid(subjectValue) {
		return nil
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(subjectValue))
	if !ok {
		return nil
	}

	roomValue, err := parser.GetCell(file, sheetName, row+1, col)
	if err != nil || !isValid(roomValue) {
		return nil
	}

	room := parser.GetCellString(roomValue)

	teacherValue, err := parser.GetCell(file, sheetName, row+3, col)
	if err != nil || !isValid(teacherValue) {
		teacherValue, err = parser.GetCell(file, sheetName, row+2, col)
		if err != nil {
			return nil
		}
	}

	if !isValid(teacherValue) {
		return nil
	}

	teacher := parser.GetCellString(teacherValue)

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacher),
		IsBlock:     true,
	}
}
