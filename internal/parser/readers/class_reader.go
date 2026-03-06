package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type ClassReader struct {
	readers []parser.Reader
}

func NewClassExtractor() *ClassReader {
	return &ClassReader{
		readers: []parser.Reader{
			SingleClassReader{},
			LectureReader{},
			BlockClassReader{},
			LargeBlockClassReader{},
		},
	}
}

func (ce *ClassReader) Extract(ws *parser.Worksheet, start types.TimeSlot, row, col int) *types.ClassSlot {
	for _, reader := range ce.readers {
		if slot, matched := reader.Read(ws, start, row, col); matched {
			return slot
		}
	}
	return nil
}
