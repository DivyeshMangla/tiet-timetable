package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/xuri/excelize/v2"
)

type LegacyReader interface {
	Read(file *excelize.File, sheetName string, row, col int) (bool, *model.ClassInfo)
}
