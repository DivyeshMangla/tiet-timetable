package model

type Timetable struct {
	Entries []TimetableEntry
}

func NewTimetable(entries []TimetableEntry) Timetable {
	return Timetable{Entries: entries}
}
