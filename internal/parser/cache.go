package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/extractor"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type Cache struct {
	batches  map[types.SheetID]map[types.BatchID]extractor.CellLocation
	daySlots map[types.SheetID][]DaySlots
}

func FromWorkbook(file *excelize.File) (*Cache, error) {
	batches := make(map[types.SheetID]map[types.BatchID]extractor.CellLocation)
	daySlots := make(map[types.SheetID][]DaySlots)

	for _, sheetName := range file.GetSheetList() {
		if isSheetHidden(file, sheetName) {
			continue
		}
		processSheet(file, sheetName, batches, daySlots)
	}

	return &Cache{batches: batches, daySlots: daySlots}, nil
}

func processSheet(
	file *excelize.File,
	sheetName string,
	batches map[types.SheetID]map[types.BatchID]extractor.CellLocation,
	daySlots map[types.SheetID][]DaySlots,
) {
	// Recover from panics in excelize or malformed sheets
	defer func() { _ = recover() }()

	sheetID := types.SheetID(sheetName)

	batchExtractor := extractor.NewBatchExtractor(file, sheetName)
	if sheetBatches, err := batchExtractor.Extract(); err == nil && len(sheetBatches) > 0 {
		batches[sheetID] = sheetBatches
	}

	firstSlotRow, firstSlotCol, found := utils.FindCellToRightOfDay(file, sheetName)
	if !found {
		return
	}

	if slots, err := BuildDaySlotsFromSheet(file, sheetName, firstSlotRow, firstSlotCol); err == nil && len(slots) > 0 {
		daySlots[sheetID] = slots
	}
}

func isSheetHidden(file *excelize.File, sheetName string) bool {
	visible, err := file.GetSheetVisible(sheetName)
	return err != nil || !visible
}
