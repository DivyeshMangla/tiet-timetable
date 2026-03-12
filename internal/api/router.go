package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
)

func SetupRoutes(timetableReg *registry.TimetableRegistry, batchReg *registry.BatchRegistry, distDir string) http.Handler {
	handler := NewHandler(timetableReg, batchReg)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/timetable/sheets", handler.GetSheets)
	mux.HandleFunc("GET /api/timetable/batches", handler.GetAllBatches)
	mux.HandleFunc("GET /api/timetable/subjects", handler.GetSubjects)
	mux.HandleFunc("GET /api/timetable/batches/{batchName}/subjects", handler.GetBatchSubjects)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches", handler.GetBatches)
	mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}", handler.GetTimetable)
	//mux.HandleFunc("GET /api/timetable/sheets/{sheetName}/batches/{batchName}/png", handler.GetTimetablePNG)
	mux.HandleFunc("POST /api/timetable/generate", handler.GetFormattedTimetablePNG)

	SetupStaticFiles(mux, distDir)

	return corsMiddleware(mux)
}

// SetupStaticFiles serves the frontend dist/ directory and falls back to
// index.html for SPA client-side routing. No-ops if distDir doesn't exist.
func SetupStaticFiles(mux *http.ServeMux, distDir string) {
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		fmt.Printf("Static files directory %q not found, skipping frontend serving\n", distDir)
		return
	}

	fs := http.FileServer(http.Dir(distDir))
	indexHTML, _ := os.ReadFile(distDir + "/index.html")

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Let API routes pass through
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the static file first
		filePath := distDir + r.URL.Path
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	fmt.Printf("Serving frontend from %s\n", distDir)
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
