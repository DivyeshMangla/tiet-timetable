package registry

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type TimetableRegistry struct {
	sheetToBatches   map[types.SheetID]map[types.BatchID]struct{}
	batchToTimetable map[types.BatchID]model.Timetable
}

func NewTimetableRegistry() *TimetableRegistry {
	return &TimetableRegistry{
		sheetToBatches:   make(map[types.SheetID]map[types.BatchID]struct{}),
		batchToTimetable: make(map[types.BatchID]model.Timetable),
	}
}

func (tr *TimetableRegistry) RegisterBatch(sheetID types.SheetID, batchID types.BatchID, timetable model.Timetable) {
	if tr.sheetToBatches[sheetID] == nil {
		tr.sheetToBatches[sheetID] = make(map[types.BatchID]struct{})
	}
	tr.sheetToBatches[sheetID][batchID] = struct{}{}
	tr.batchToTimetable[batchID] = timetable
}

func (tr *TimetableRegistry) GetTimetable(batchID types.BatchID) (model.Timetable, bool) {
	timetable, ok := tr.batchToTimetable[batchID]
	return timetable, ok
}

func (tr *TimetableRegistry) GetBatches(sheetID types.SheetID) map[types.BatchID]struct{} {
	return tr.sheetToBatches[sheetID]
}

func (tr *TimetableRegistry) SheetIDs() []types.SheetID {
	ids := make([]types.SheetID, 0, len(tr.sheetToBatches))
	for id := range tr.sheetToBatches {
		ids = append(ids, id)
	}
	return ids
}

func (tr *TimetableRegistry) BatchIDs(sheetID types.SheetID) []types.BatchID {
	batches, ok := tr.sheetToBatches[sheetID]
	if !ok {
		return nil
	}
	ids := make([]types.BatchID, 0, len(batches))
	for id := range batches {
		ids = append(ids, id)
	}
	return ids
}
