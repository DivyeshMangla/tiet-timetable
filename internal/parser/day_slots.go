package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/xuri/excelize/v2"
)

type DaySlots struct {
	Day   model.Day
	Slots map[model.TimeSlot]CellLocation
}

func BuildDaySlotsFromSheet(file *excelize.File, sheetName string, firstSlotRow, firstSlotCol int) ([]DaySlots, error) {
	if firstSlotRow < 0 || firstSlotCol < 0 {
		return nil, nil
	}

	factory := newDaySlotsFactory()
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	for row := firstSlotRow; row < len(rows) && !factory.isComplete(); row++ {
		value, err := GetCell(file, sheetName, row, firstSlotCol)
		if err != nil {
			continue
		}

		slotNumber, ok := ParseSlotNumber(value)
		if !ok {
			continue
		}

		factory.processSlot(slotNumber, CellLocation{
			Row: row,
			Col: firstSlotCol,
		})
	}

	return factory.build(), nil
}
