package parser

import (
	"strings"

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

type daySlotsFactory struct {
	result       []DaySlots
	dayIndex     int
	currentSlots map[model.TimeSlot]CellLocation
}

var weekdays = []model.Day{
	model.MON, model.TUE, model.WED, model.THU, model.FRI,
}

func newDaySlotsFactory() *daySlotsFactory {
	return &daySlotsFactory{
		result:       make([]DaySlots, 0, len(weekdays)),
		currentSlots: make(map[model.TimeSlot]CellLocation),
	}
}

func (f *daySlotsFactory) processSlot(slotNumber int, cell CellLocation) {
	if slotNumber == 1 && len(f.currentSlots) > 0 {
		f.finalizeCurrentDay()
	}

	slot, err := model.FromNumber(slotNumber)
	if err != nil {
		return
	}

	f.currentSlots[slot] = cell
}

func (f *daySlotsFactory) finalizeCurrentDay() {
	if f.dayIndex >= len(weekdays) {
		return
	}

	f.result = append(f.result, DaySlots{
		Day:   weekdays[f.dayIndex],
		Slots: f.currentSlots,
	})

	f.dayIndex++
	f.currentSlots = make(map[model.TimeSlot]CellLocation)
}

func (f *daySlotsFactory) isComplete() bool {
	return f.dayIndex >= len(weekdays)
}

func (f *daySlotsFactory) build() []DaySlots {
	if len(f.currentSlots) > 0 && f.dayIndex < len(weekdays) {
		f.result = append(f.result, DaySlots{
			Day:   weekdays[f.dayIndex],
			Slots: f.currentSlots,
		})
	}
	return f.result
}

const maxRowsToScan = 300

func BuildDaySlotsFromSheet(file *excelize.File, sheetName string, firstSlotRow, firstSlotCol int) ([]DaySlots, error) {
	if file == nil || firstSlotRow < 0 || firstSlotCol < 0 {
		return nil, nil
	}

	factory := newDaySlotsFactory()
	maxRows := firstSlotRow + maxRowsToScan
	lastSlotNumber := -1

	for row := firstSlotRow; row <= maxRows && !factory.isComplete(); row++ {
		value, _ := utils.GetCell(file, sheetName, row, firstSlotCol)

		if strings.TrimSpace(value) == "" {
			continue
		}

		slotNumber, ok := utils.ParseSlotNumber(value)
		if !ok {
			continue
		}

		if slotNumber == lastSlotNumber {
			continue
		}
		lastSlotNumber = slotNumber

		factory.processSlot(slotNumber, CellLocation{
			Row: row,
			Col: firstSlotCol,
		})
	}

	return factory.build(), nil
}
