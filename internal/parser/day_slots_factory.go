package parser

import "github.com/DivyeshMangla/tiet-timetable/internal/model"

type daySlotsFactory struct {
	result          []DaySlots
	days            []model.Day
	currentDayIndex int
	currentDaySlots map[model.TimeSlot]CellLocation
}

func newDaySlotsFactory() *daySlotsFactory {
	return &daySlotsFactory{
		result:          make([]DaySlots, 0),
		days:            []model.Day{model.MON, model.TUE, model.WED, model.THU, model.FRI, model.SAT},
		currentDayIndex: 0,
		currentDaySlots: make(map[model.TimeSlot]CellLocation),
	}
}

func (f *daySlotsFactory) processSlot(slotNumber int, cell CellLocation) {
	if f.isDayBoundary(slotNumber) {
		f.finalizeCurrentDay()
	}

	slot, err := model.FromNumber(slotNumber)
	if err != nil {
		return
	}

	f.currentDaySlots[slot] = cell
}

func (f *daySlotsFactory) isDayBoundary(slotNumber int) bool {
	return slotNumber == 1 && len(f.currentDaySlots) > 0
}

func (f *daySlotsFactory) finalizeCurrentDay() {
	if f.currentDayIndex >= len(f.days) {
		return
	}

	f.result = append(f.result, DaySlots{
		Day:   f.days[f.currentDayIndex],
		Slots: f.currentDaySlots,
	})

	f.currentDayIndex++
	f.currentDaySlots = make(map[model.TimeSlot]CellLocation)
}

func (f *daySlotsFactory) isComplete() bool {
	return f.currentDayIndex >= len(f.days)
}

func (f *daySlotsFactory) build() []DaySlots {
	if len(f.currentDaySlots) > 0 && f.currentDayIndex < len(f.days) {
		f.result = append(f.result, DaySlots{
			Day:   f.days[f.currentDayIndex],
			Slots: f.currentDaySlots,
		})
	}

	return f.result
}
