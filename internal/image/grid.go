package image

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
)

type GridBounds struct {
	YBounds []Bounds
	XBounds []Bounds
}

type Bounds struct {
	Start int
	End   int
}

var ScheduleGrid = GridBounds{
	YBounds: []Bounds{
		{Start: 736, End: 915},
		{Start: 944, End: 1115},
		{Start: 1144, End: 1323},
		{Start: 1352, End: 1523},
		{Start: 1552, End: 1731},
		{Start: 1760, End: 1931},
		{Start: 1960, End: 2131},
		{Start: 2160, End: 2339},
		{Start: 2364, End: 2543},
		{Start: 2568, End: 2747},
		{Start: 2776, End: 2947},
	},

	XBounds: []Bounds{
		{Start: 928, End: 1631},
		{Start: 1664, End: 2371},
		{Start: 2400, End: 3107},
		{Start: 3136, End: 3847},
		{Start: 3876, End: 4579},
	},
}

type Cell struct {
	X      int
	Y      int
	Width  int
	Height int
}

func GetCell(timeSlot model.TimeSlot, day model.Day) Cell {
	if int(timeSlot) >= len(ScheduleGrid.YBounds) || int(day) >= len(ScheduleGrid.XBounds) {
		return Cell{}
	}

	yBound := ScheduleGrid.YBounds[timeSlot]
	xBound := ScheduleGrid.XBounds[day]

	return Cell{
		X:      xBound.Start,
		Y:      yBound.Start,
		Width:  xBound.End - xBound.Start,
		Height: yBound.End - yBound.Start,
	}
}
