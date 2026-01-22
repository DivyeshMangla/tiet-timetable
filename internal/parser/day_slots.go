package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/xuri/excelize/v2"
)

type DaySlots struct {
	Day   model.Day
	Slots map[model.TimeSlot]CellLocation
}

type CellLocation struct {
	Row int
	Col int
}

func BuildDaySlotsFromSheet(file *excelize.File, sheetName string, firstSlotRow, firstSlotCol int) ([]DaySlots, error) {
	if firstSlotRow < 0 || firstSlotCol < 0 {
		return nil, nil
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	factory := newDaySlotsFactory()
	for row := firstSlotRow; row < len(rows) && !factory.isComplete(); row++ {
		value, err := utils.GetCell(file, sheetName, row, firstSlotCol)
		if err != nil {
			continue
		}

		slotNumber, ok := utils.ParseSlotNumber(value)
		if !ok {
			continue
		}

		factory.processSlot(slotNumber, CellLocation{Row: row, Col: firstSlotCol})
	}

	return factory.build(), nil
}
