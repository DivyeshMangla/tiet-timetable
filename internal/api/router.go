package api

import (
	"net/http"

	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
)

func SetupRoutes(reg *registry.TimetableRegistry) http.Handler {
	handler := NewHandler(reg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/timetable/sheets", handler.GetSheets)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches", handler.GetBatches)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}", handler.GetTimetable)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}/png", handler.GetTimetablePNG)

	return mux
}
