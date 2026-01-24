package readers

import (
	"regexp"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

const minHorizontalMergeWidth = 3

var teacherPattern = regexp.MustCompile(`^[A-Za-z. ]+$`)

type LargeClassReader struct{}

func (r *LargeClassReader) Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo) {
	region, found := utils.GetHorizontalMergedRegion(file, sheetName, row, col)
	if !found || !isWideEnough(region) {
		return false, nil
	}

	classCodeValue, err := utils.GetCell(file, sheetName, row, region.StartCol)
	if err != nil {
		return false, nil
	}

	parsed, ok := utils.ParseCode(utils.GetCellString(utils.FirstLine(classCodeValue)))
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
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

	teacher := utils.GetCellString(utils.FirstLine(teacherValue))
	if !isValidTeacher(teacher) {
		return false, nil
	}

	room := utils.GetCellString(utils.FirstLine(roomValue))
	return true, &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(room),
		Teacher:     types.Teacher(teacher),
		IsBlock:     false,
	}
}

func isWideEnough(region utils.MergedRegion) bool {
	return (region.EndCol - region.StartCol) >= minHorizontalMergeWidth
}

func isValidTeacher(teacher string) bool {
	trimmed := strings.TrimSpace(teacher)
	return trimmed != "" && teacherPattern.MatchString(trimmed)
}
