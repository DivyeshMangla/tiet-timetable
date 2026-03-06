package registry

import "github.com/DivyeshMangla/tiet-timetable/internal/types"

type BatchRegistry struct {
	sheetToBatches map[types.SheetID][]types.BatchID
}

func NewBatchRegistry() *BatchRegistry {
	return &BatchRegistry{
		sheetToBatches: make(map[types.SheetID][]types.BatchID),
	}
}

// AddBatch adds a batch to a sheet
func (r *BatchRegistry) AddBatch(sheet types.SheetID, batch types.BatchID) {
	r.sheetToBatches[sheet] = append(r.sheetToBatches[sheet], batch)
}

// TotalCount returns total number of batches across all sheets
func (r *BatchRegistry) TotalCount() int {
	total := 0
	for _, batches := range r.sheetToBatches {
		total += len(batches)
	}
	return total
}

// AllBatches returns all batches across all sheets (flattened)
func (r *BatchRegistry) AllBatches() []types.BatchID {
	total := r.TotalCount()
	all := make([]types.BatchID, 0, total)

	for _, batches := range r.sheetToBatches {
		all = append(all, batches...)
	}

	return all
}
