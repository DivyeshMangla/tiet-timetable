package reader

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type Reader interface {
	Read(ws *excel.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool)
}
