package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type Reader interface {
	Read(ws *Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool)
}
