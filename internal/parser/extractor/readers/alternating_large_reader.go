package readers

import (
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type AlternatingLargeClassReader struct{}

func (r *AlternatingLargeClassReader) Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo) {
	region, found := utils.GetHorizontalMergedRegion(file, sheetName, row, col)
	if !found || !isWideEnough(region) {
		return false, nil
	}

	classCodeValue, err := utils.GetCell(file, sheetName, row, region.StartCol)
	if err != nil {
		return false, nil
	}

	roomValue, err := utils.GetCell(file, sheetName, row+1, region.StartCol)
	if err != nil {
		return false, nil
	}

	teacherValue, err := utils.GetCell(file, sheetName, row+1, region.EndCol)
	if err != nil {
		return false, nil
	}

	subjectStr := utils.GetCellString(utils.FirstLine(classCodeValue))
	_ = utils.GetCellString(utils.FirstLine(roomValue))
	teacherStr := utils.GetCellString(utils.FirstLine(teacherValue))

	if !strings.Contains(subjectStr, "/") {
		return false, nil
	}

	subjects := splitAndTrim(subjectStr, "/")
	if len(subjects) < 2 {
		return false, nil
	}

	parsed, ok := utils.ParseCode(subjects[0])
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
		return false, nil
	}

	return true, &model.ClassInfo{
		SubjectCode: types.SubjectCode("Elective"),
		ClassType:   parsed.ClassType,
		Room:        types.Room("Elective"),
		Teacher:     types.Teacher(teacherStr),
		IsBlock:     false,
	}
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
