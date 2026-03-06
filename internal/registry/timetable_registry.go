package registry

import "github.com/DivyeshMangla/tiet-timetable/internal/types"

type TimetableRegistry struct {
	batchToTimetable map[types.BatchID]*types.Timetable
}

func NewTimetableRegistry() *TimetableRegistry {
	return &TimetableRegistry{
		batchToTimetable: make(map[types.BatchID]*types.Timetable),
	}
}

func (r *TimetableRegistry) AddTimetable(batch types.BatchID, timetable *types.Timetable) {
	r.batchToTimetable[batch] = timetable
}

func (r *TimetableRegistry) Get(batch types.BatchID) (*types.Timetable, bool) {
	timetable, ok := r.batchToTimetable[batch]
	return timetable, ok
}

func (r *TimetableRegistry) TotalCount() int {
	return len(r.batchToTimetable)
}

func (r *TimetableRegistry) AllTimetables() []*types.Timetable {
	all := make([]*types.Timetable, 0, len(r.batchToTimetable))
	for _, timetable := range r.batchToTimetable {
		all = append(all, timetable)
	}
	return all
}

func (r *TimetableRegistry) AllBatches() []types.BatchID {
	batches := make([]types.BatchID, 0, len(r.batchToTimetable))
	for batch := range r.batchToTimetable {
		batches = append(batches, batch)
	}
	return batches
}

func (r *TimetableRegistry) AllUniqueSubjects() []types.SubjectCode {
	seen := make(map[types.SubjectCode]struct{})

	for _, timetable := range r.batchToTimetable {
		for _, code := range timetable.AllUniqueSubjects() {
			seen[code] = struct{}{}
		}
	}

	codes := make([]types.SubjectCode, 0, len(seen))
	for code := range seen {
		codes = append(codes, code)
	}

	return codes
}
