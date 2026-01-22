package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type BlockClassReader struct{}

func (r *BlockClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	subjectValue, err := utils.GetCell(file, sheetName, row, col)
	if err != nil || !utils.IsValid(subjectValue) {
		return false
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(subjectValue))
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
		return false
	}

	roomValue, err := utils.GetCell(file, sheetName, row+1, col)
	if err != nil || !utils.IsValid(roomValue) {
		return false
	}

	row2Value, _ := utils.GetCell(file, sheetName, row+2, col)
	row3Value, _ := utils.GetCell(file, sheetName, row+3, col)

	if utils.IsValid(row2Value) && utils.IsValid(row3Value) {
		return true
	}

	return utils.IsValid(row2Value) && !utils.IsValid(row3Value)
}

func (r *BlockClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	subjectValue, err := utils.GetCell(file, sheetName, row, col)
	if err != nil || !utils.IsValid(subjectValue) {
		return nil
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(subjectValue))
	if !ok {
		return nil
	}

	roomValue, err := utils.GetCell(file, sheetName, row+1, col)
	if err != nil || !utils.IsValid(roomValue) {
		return nil
	}

	room := utils.GetCellString(roomValue)

	teacherValue, err := utils.GetCell(file, sheetName, row+3, col)
	if err != nil || !utils.IsValid(teacherValue) {
		teacherValue, err = utils.GetCell(file, sheetName, row+2, col)
		if err != nil {
			return nil
		}
	}

	if !utils.IsValid(teacherValue) {
		return nil
	}

	teacher := utils.GetCellString(teacherValue)

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacher),
		IsBlock:     true,
	}
}
