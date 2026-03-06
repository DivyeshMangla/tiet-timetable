package parser

import "github.com/xuri/excelize/v2"

type Parser struct {
	workbook *excelize.File
	layout   *WorkbookLayout
}

func NewParser(file *excelize.File, layout *WorkbookLayout) *Parser {
	return &Parser{
		workbook: file,
		layout:   layout,
	}
}
