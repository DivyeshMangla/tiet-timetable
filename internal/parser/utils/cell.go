package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var subjectCodePattern = regexp.MustCompile(`^[A-Z]{3}\d{3}$`)

func GetCell(file *excelize.File, sheetName string, row, col int) (string, error) {
	cellRef := ToCellRef(row, col)
	return file.GetCellValue(sheetName, cellRef)
}

func GetCellString(value string) string {
	return strings.TrimSpace(value)
}

func ParseSlotNumber(value string) (int, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false
	}

	num, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, false
	}

	return num, true
}

func FindCellInFirstColumn(file *excelize.File, sheetName string, searchText string) (int, int, bool) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return 0, 0, false
	}

	searchLower := strings.ToLower(strings.TrimSpace(searchText))

	for row := 0; row < len(rows); row++ {
		if len(rows[row]) == 0 {
			continue
		}

		value := strings.TrimSpace(rows[row][0])
		if strings.ToLower(value) == searchLower {
			return row, 0, true
		}
	}

	return 0, 0, false
}

func IsSubjectCode(code string) bool {
	return subjectCodePattern.MatchString(strings.TrimSpace(code))
}

func FindCellToRightOfDay(file *excelize.File, sheetName string) (int, int, bool) {
	dayRow, dayCol, found := FindCellInFirstColumn(file, sheetName, "day")
	if !found {
		return 0, 0, false
	}

	return dayRow, dayCol + 1, true
}

func ToCellRef(row, col int) string {
	colName := ""
	colNum := col + 1
	for colNum > 0 {
		colNum--
		colName = string(rune('A'+colNum%26)) + colName
		colNum /= 26
	}
	return colName + strconv.Itoa(row+1)
}
