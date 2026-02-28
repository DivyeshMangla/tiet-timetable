package api

import (
	"net/http"

	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
)

func SetupRoutes(reg *registry.TimetableRegistry) http.Handler {
	handler := NewHandler(reg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/timetable/sheets", handler.GetSheets)
	mux.HandleFunc("GET /api/timetable/batches", handler.GetAllBatches)
	mux.HandleFunc("GET /api/timetable/subjects", handler.GetSubjects)
	mux.HandleFunc("GET /api/timetable/batches/{batchName}/subjects", handler.GetBatchSubjects)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches", handler.GetBatches)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}", handler.GetTimetable)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}/png", handler.GetTimetablePNG)
	mux.HandleFunc("POST /api/timetable/generate", handler.GetFormattedTimetablePNG)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
