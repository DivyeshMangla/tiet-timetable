package parser

import "strconv"

type CellLocation struct {
	Row int
	Col int
}

func (c CellLocation) ToCellRef() string {
	return ToCellRef(c.Row, c.Col)
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
