package types

type (
	Class struct {
		SubjectCode SubjectCode
		Room        Room
		Teacher     Teacher
	}

	ClassSlot struct {
		Start   TimeSlot
		End     TimeSlot
		Classes []Class
	}

	Timetable struct {
		Batch BatchID
		Days  map[Day][]ClassSlot
	}
)

func (t *Timetable) AllUniqueSubjects() []SubjectCode {
	seen := make(map[SubjectCode]struct{})

	for _, slots := range t.Days {
		for _, slot := range slots {
			for _, class := range slot.Classes {
				seen[class.SubjectCode] = struct{}{}
			}
		}
	}

	codes := make([]SubjectCode, 0, len(seen))
	for code := range seen {
		codes = append(codes, code)
	}

	return codes
}
