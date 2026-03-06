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

// SheetIDs returns all registered sheet IDs
func (r *BatchRegistry) SheetIDs() []types.SheetID {
	sheets := make([]types.SheetID, 0, len(r.sheetToBatches))
	for sheet := range r.sheetToBatches {
		sheets = append(sheets, sheet)
	}
	return sheets
}

// BatchesBySheet returns all batches for a given sheet
func (r *BatchRegistry) BatchesBySheet(sheet types.SheetID) []types.BatchID {
	return r.sheetToBatches[sheet]
}
