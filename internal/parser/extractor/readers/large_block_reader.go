package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type LargeBlockClassReader struct{}

func (r *LargeBlockClassReader) Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo) {
	region, found := utils.GetHorizontalMergedRegion(file, sheetName, row, col)
	if !found || !isWideEnough(region) {
		return false, nil
	}

	startRow := region.StartRow
	startCol := region.StartCol

	classCodeValue, err := utils.GetCell(file, sheetName, startRow, startCol)
	if err != nil {
		return false, nil
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(utils.FirstLine(classCodeValue)))
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
		return false, nil
	}

	roomValue, err := utils.GetCell(file, sheetName, startRow+1, startCol)
	if err != nil {
		return false, nil
	}

	hasContinuation := utils.GetCellString(utils.FirstLine(
		mustGetCell(file, sheetName, startRow+2, startCol))) != "" ||
		utils.GetCellString(utils.FirstLine(
			mustGetCell(file, sheetName, startRow+3, startCol))) != ""

	if !hasContinuation {
		return false, nil
	}

	// Prefer teacher from row+3, fallback to row+2
	teacherValue, err := utils.GetCell(file, sheetName, startRow+3, startCol)
	if err != nil {
		return false, nil
	}
	teacher := utils.GetCellString(utils.FirstLine(teacherValue))

	if teacher == "" {
		teacherValue, err = utils.GetCell(file, sheetName, startRow+2, startCol)
		if err != nil {
			return false, nil
		}
		teacher = utils.GetCellString(utils.FirstLine(teacherValue))
	}

	if !isValidTeacher(teacher) {
		return false, nil
	}

	room := utils.GetCellString(utils.FirstLine(roomValue))
	return true, &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacher),
		IsBlock:     true,
	}
}

func mustGetCell(file *excelize.File, sheetName string, row, col int) string {
	val, _ := utils.GetCell(file, sheetName, row, col)
	return val
}
