package utils

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

func GetCell(file *excelize.File, sheetName string, row, col int) (string, error) {
	cellRef := ToCellRef(row, col)
	return file.GetCellValue(sheetName, cellRef)
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
