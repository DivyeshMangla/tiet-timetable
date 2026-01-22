package readers

import (
	"regexp"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

var teacherPattern = regexp.MustCompile(`^[A-Za-z. ]+$`)

type LargeClassReader struct{}

func (r *LargeClassReader) Matches(file *excelize.File, sheetName string, row, col int) bool {
	region, found := parser.GetHorizontalMergedRegion(file, sheetName, row, col)
	if !found || !isWideEnough(region) {
		return false
	}

	subjectValue, err := parser.GetCell(file, sheetName, row, region.StartCol)
	if err != nil {
		return false
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(subjectValue))
	if !ok {
		return false
	}

	return parser.IsSubjectCode(string(parsed.Code))
}

func (r *LargeClassReader) Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	region, found := parser.GetHorizontalMergedRegion(file, sheetName, row, col)
	if !found {
		return nil
	}

	classCodeValue, err := parser.GetCell(file, sheetName, row, region.StartCol)
	if err != nil {
		return nil
	}

	roomValue, err := parser.GetCell(file, sheetName, row+1, region.StartCol)
	if err != nil {
		return nil
	}

	teacherValue, err := parser.GetCell(file, sheetName, row+1, region.EndCol)
	if err != nil {
		return nil
	}

	parsed, ok := parser.ParseCode(parser.GetCellString(classCodeValue))
	if !ok {
		return nil
	}

	teacher := parser.GetCellString(teacherValue)
	if !isValidTeacher(teacher) {
		return nil
	}

	return &model.ClassInfo{
		SubjectCode: parsed.Code,
		ClassType:   parsed.ClassType,
		Room:        types.Room(parser.GetCellString(roomValue)),
		Teacher:     types.Teacher(teacher),
		IsBlock:     false,
	}
}

func isWideEnough(region parser.MergedRegion) bool {
	return (region.EndCol - region.StartCol) > 2
}

func isValidTeacher(teacher string) bool {
	trimmed := strings.TrimSpace(teacher)
	return trimmed != "" && teacherPattern.MatchString(trimmed)
}
