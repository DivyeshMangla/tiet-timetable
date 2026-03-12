package types

type (
	RenderInfo struct {
		Start     TimeSlot
		End       TimeSlot
		ClassType ClassType
		Text      string
	}

	RenderableTimetable struct {
		Batch BatchID
		Days  map[Day][]RenderInfo
	}
)

func (ri RenderInfo) IsBlock() bool {
	return ri.End-ri.Start >= 1
}
