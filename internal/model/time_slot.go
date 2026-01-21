package model

import (
	"fmt"
	"time"
)

type TimeSlot int

const (
	T1 TimeSlot = iota
	T2
	T3
	T4
	T5
	T6
	T7
	T8
	T9
	T10
	T11
)

var timeSlotTimes = []struct {
	start string
	end   string
}{
	{"08:00", "08:50"},
	{"08:50", "09:40"},
	{"09:40", "10:30"},
	{"10:30", "11:20"},
	{"11:20", "12:10"},
	{"12:10", "13:00"},
	{"13:00", "13:50"},
	{"13:50", "14:40"},
	{"14:40", "15:30"},
	{"15:30", "16:20"},
	{"16:20", "17:10"},
}

const timeFormat = "15:04"

func (ts TimeSlot) Start() time.Time {
	if ts < 0 || int(ts) >= len(timeSlotTimes) {
		return time.Time{}
	}
	t, _ := time.Parse(timeFormat, timeSlotTimes[ts].start)
	return t
}

func (ts TimeSlot) End() time.Time {
	if ts < 0 || int(ts) >= len(timeSlotTimes) {
		return time.Time{}
	}
	t, _ := time.Parse(timeFormat, timeSlotTimes[ts].end)
	return t
}

func FromNumber(number int) (TimeSlot, error) {
	if number < 1 || number > len(timeSlotTimes) {
		return 0, fmt.Errorf("invalid time slot number: %d", number)
	}

	return TimeSlot(number - 1), nil
}

func (ts TimeSlot) Label() string {
	if ts < 0 || int(ts) >= len(timeSlotTimes) {
		return ""
	}

	return fmt.Sprintf("%s–%s", timeSlotTimes[ts].start, timeSlotTimes[ts].end)
}
