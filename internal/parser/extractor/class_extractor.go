package extractor

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/extractor/readers"
	"github.com/xuri/excelize/v2"
)

type ClassExtractor struct {
	readers []readers.Reader
}

func NewClassExtractor() *ClassExtractor {
	return &ClassExtractor{
		readers: []readers.Reader{
			&readers.SingleClassReader{},
			&readers.LargeClassReader{},
			&readers.BlockClassReader{},
		},
	}
}

func (ce *ClassExtractor) Extract(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	for _, reader := range ce.readers {
		if reader.Matches(file, sheetName, row, col) {
			return reader.Read(file, sheetName, row, col)
		}
	}

	return nil
}
