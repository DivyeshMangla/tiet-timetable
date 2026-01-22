package api

import (
	"encoding/json"
	"net/http"

	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type Handler struct {
	registry *registry.TimetableRegistry
}

func NewHandler(reg *registry.TimetableRegistry) *Handler {
	return &Handler{registry: reg}
}

func (h *Handler) GetSheets(w http.ResponseWriter, r *http.Request) {
	sheetIDs := h.registry.SheetIDs()
	sheetNames := make([]string, 0, len(sheetIDs))

	for _, sheetID := range sheetIDs {
		sheetNames = append(sheetNames, string(sheetID))
	}

	writeJSON(w, http.StatusOK, sheetNames)
}

func (h *Handler) GetBatches(w http.ResponseWriter, r *http.Request) {
	sheetName := r.PathValue("sheetName")
	if sheetName == "" {
		writeError(w, http.StatusBadRequest, "sheet name is required")
		return
	}

	sheetID := types.SheetID(sheetName)
	batches := h.registry.GetBatches(sheetID)
	if batches == nil {
		writeError(w, http.StatusNotFound, "sheet not found")
		return
	}

	batchNames := make([]string, 0, len(batches))
	for batchID := range batches {
		batchNames = append(batchNames, string(batchID))
	}

	writeJSON(w, http.StatusOK, batchNames)
}

func (h *Handler) GetTimetable(w http.ResponseWriter, r *http.Request) {
	sheetName := r.PathValue("sheetName")
	batchName := r.PathValue("batchName")

	if sheetName == "" || batchName == "" {
		writeError(w, http.StatusBadRequest, "sheet name and batch name are required")
		return
	}

	batchID := types.BatchID(batchName)
	timetable, ok := h.registry.GetTimetable(batchID)
	if !ok {
		writeError(w, http.StatusNotFound, "timetable not found")
		return
	}

	writeJSON(w, http.StatusOK, timetable.Entries)
}

func (h *Handler) GetTimetablePNG(w http.ResponseWriter, r *http.Request) {
	sheetName := r.PathValue("sheetName")
	batchName := r.PathValue("batchName")

	if sheetName == "" || batchName == "" {
		writeError(w, http.StatusBadRequest, "sheet name and batch name are required")
		return
	}

	batchID := types.BatchID(batchName)
	_, ok := h.registry.GetTimetable(batchID)
	if !ok {
		writeError(w, http.StatusNotFound, "timetable not found")
		return
	}

	// TODO: Implement PNG generation
	writeError(w, http.StatusNotImplemented, "PNG generation not yet implemented")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
