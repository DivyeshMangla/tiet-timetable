package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/xuri/excelize/v2"
)

type Reader interface {
	Matches(file *excelize.File, sheetName string, row, col int) bool
	Read(file *excelize.File, sheetName string, row, col int) *model.ClassInfo
}
