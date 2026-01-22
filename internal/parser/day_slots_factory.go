package parser

import "github.com/DivyeshMangla/tiet-timetable/internal/model"

var weekdays = []model.Day{
	model.MON, model.TUE, model.WED, model.THU, model.FRI, model.SAT,
}

type daySlotsFactory struct {
	result       []DaySlots
	dayIndex     int
	currentSlots map[model.TimeSlot]CellLocation
}

func newDaySlotsFactory() *daySlotsFactory {
	return &daySlotsFactory{
		result:       make([]DaySlots, 0, len(weekdays)),
		dayIndex:     0,
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
