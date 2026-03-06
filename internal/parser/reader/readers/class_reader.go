package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/reader"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type ClassReader struct {
	readers []reader.Reader
}

func NewClassExtractor() *ClassReader {
	return &ClassReader{
		readers: []reader.Reader{
			SingleClassReader{},
			LectureReader{},
			BlockClassReader{},
			LargeBlockClassReader{},
		},
	}
}

func (ce *ClassReader) Extract(ws *excel.Worksheet, start types.TimeSlot, row, col int) *types.ClassSlot {
	for _, reader := range ce.readers {
		if slot, matched := reader.Read(ws, start, row, col); matched {
			return slot
		}
	}
	return nil
}
