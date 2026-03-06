package types

type (
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
)

func (c ClassSlot) IsBlock() bool {
	return c.End-c.Start >= 1
}
