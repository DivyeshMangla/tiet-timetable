package internal

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

func Bootstrap(file *excelize.File) error {
	layout, err := parser.NewWorkbookLayoutBuilder(file).Build()
	if err != nil {
		return err
	}

	// Make & Populate BatchRegistry
	batchRegistry := registry.NewBatchRegistry()
	for sheet, sheetLayout := range layout.Sheets {
		for batchID := range sheetLayout.BatchCells {
			batchRegistry.AddBatch(types.SheetID(sheet), batchID)
		}
	}

	return nil
}
