package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/xuri/excelize/v2"
)

func getValidCell(file *excelize.File, sheet string, row, col int) string {
	val, err := utils.GetCell(file, sheet, row, col)
	if err != nil || !utils.IsValid(val) {
		return ""
	}
	return val
}

func parseSubject(value string) (utils.ParsedCode, bool) {
	line := utils.GetCellString(utils.FirstLine(value))
	parsed, ok := utils.ParseCode(line)
	if !ok || !utils.IsSubjectCode(string(parsed.Code)) {
		return parsed, false
	}
	return parsed, true
}

func cleanCell(value string) string {
	return utils.GetCellString(utils.FirstLine(value))
}
