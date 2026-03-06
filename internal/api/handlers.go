package api

import (
	"encoding/json"
	"fmt"
	"github.com/DivyeshMangla/tiet-timetable/internal/image"
	"net/http"
	"os"

	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type Handler struct {
	timetableReg *registry.TimetableRegistry
	batchReg     *registry.BatchRegistry
}

func NewHandler(timetableReg *registry.TimetableRegistry, batchReg *registry.BatchRegistry) *Handler {
	return &Handler{timetableReg: timetableReg, batchReg: batchReg}
}

func (h *Handler) GetSheets(w http.ResponseWriter, r *http.Request) {
	sheetIDs := h.batchReg.SheetIDs()
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
	batches := h.batchReg.BatchesBySheet(sheetID)
	if batches == nil {
		writeError(w, http.StatusNotFound, "sheet not found")
		return
	}

	batchNames := make([]string, 0, len(batches))
	for batchID := range batches {
		batchNames = append(batchNames, string(rune(batchID)))
	}

	writeJSON(w, http.StatusOK, batchNames)
}

func (h *Handler) GetAllBatches(w http.ResponseWriter, r *http.Request) {
	batches := h.timetableReg.AllBatches()
	batchNames := make([]string, 0, len(batches))
	for _, batchID := range batches {
		batchNames = append(batchNames, string(batchID))
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"batches": batchNames})
}

func (h *Handler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	subjects := h.timetableReg.AllUniqueSubjects()
	codes := make([]string, 0, len(subjects))
	for _, code := range subjects {
		codes = append(codes, string(code))
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"subjects": codes})
}

func (h *Handler) GetBatchSubjects(w http.ResponseWriter, r *http.Request) {
	batchName := r.PathValue("batchName")
	if batchName == "" {
		writeError(w, http.StatusBadRequest, "batch name is required")
		return
	}

	timetable, ok := h.timetableReg.Get(types.BatchID(batchName))
	if !ok {
		writeError(w, http.StatusNotFound, "batch not found")
		return
	}
	subjects := timetable.AllUniqueSubjects()

	codes := make([]string, 0, len(subjects))
	for _, code := range subjects {
		codes = append(codes, string(code))
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"subjects": codes})
}

func (h *Handler) GetTimetable(w http.ResponseWriter, r *http.Request) {
	sheetName := r.PathValue("sheetName")
	batchName := r.PathValue("batchName")

	if sheetName == "" || batchName == "" {
		writeError(w, http.StatusBadRequest, "sheet name and batch name are required")
		return
	}

	batchID := types.BatchID(batchName)
	timetable, ok := h.timetableReg.Get(batchID)
	if !ok {
		writeError(w, http.StatusNotFound, "timetable not found")
		return
	}

	writeJSON(w, http.StatusOK, timetable)
}

func (h *Handler) GetFormattedTimetablePNG(w http.ResponseWriter, r *http.Request) {
	var req FormattedTimetableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Batch == "" || len(req.Subjects) == 0 {
		writeError(w, http.StatusBadRequest, "batch and subjects are required")
		return
	}

	batchID := types.BatchID(req.Batch)
	timetable, ok := h.timetableReg.Get(batchID)
	if !ok {
		writeError(w, http.StatusNotFound, "timetable not found")
		return
	}

	drawer, err := image.NewTimetableDrawer()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create timetable drawer")
		return
	}

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("formatted_%s_*.png", req.Batch))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create temp file")
		return
	}
	outputPath := tmpFile.Name()
	tmpFile.Close()

	defer os.Remove(outputPath)

	err = drawer.DrawTimetable(BuildRenderableTimetable(timetable, req.Subjects), outputPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate timetable image")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, outputPath)
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
