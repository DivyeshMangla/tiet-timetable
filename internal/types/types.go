package types

type (
	ClassType   int
	TimeSlot    int
	SheetID     string
	BatchID     string
	SubjectCode string
	Room        string
	Teacher     string

	CellLocation struct {
		Row int
		Col int
	}

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
)
