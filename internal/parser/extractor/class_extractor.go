package extractor

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/extractor/readers"
	"github.com/xuri/excelize/v2"
)

// ClassExtractor attempts to extract class information from a cell using
// multiple reader strategies. Readers are tried in order until one matches.
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
