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
			&readers.AlternatingLargeClassReader{},
		},
	}
}

func (ce *ClassExtractor) Extract(file *excelize.File, sheetName string, row, col int) *model.ClassInfo {
	for _, reader := range ce.readers {
		if matched, classInfo := reader.Read(file, sheetName, row, col); matched {
			return classInfo
		}
	}
	return nil
}
