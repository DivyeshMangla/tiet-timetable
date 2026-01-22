package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type SingleClassReader struct{}

func (r *SingleClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	subjectValue, err := utils.GetCell(file, sheetName, row, col)
	if err != nil {
		return false
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(subjectValue))
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
		return false
	}

	roomValue, err := utils.GetCell(file, sheetName, row+1, col)
	if err != nil {
		return false
	}

	teacherValue, err := utils.GetCell(file, sheetName, row+1, col+1)
	if err != nil {
		return false
	}

	if !utils.IsValid(roomValue) || !utils.IsValid(teacherValue) {
		return false
	}

	return true
}

func (r *SingleClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subjectValue, err := utils.GetCell(file, sheetName, row, col)
	if err != nil {
		return nil
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(subjectValue))
	if !ok {
		return nil
	}

	roomValue, _ := utils.GetCell(file, sheetName, row+1, col)
	teacherValue, _ := utils.GetCell(file, sheetName, row+1, col+1)

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(utils.GetCellString(roomValue)),
		Teacher:     types.Teacher(utils.GetCellString(teacherValue)),
		IsBlock:     false,
	}
}
