package internal

import (
	"fmt"
	"github.com/DivyeshMangla/tiet-timetable/internal/api"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

func Bootstrap(file *excelize.File) error {
	layout, err := parser.NewWorkbookLayoutBuilder(file).Build()
	if err != nil {
		return err
	}

	// Make & Populate BatchRegistry
	batchRegistry := registry.NewBatchRegistry()
	for sheet, sheetLayout := range layout.Sheets {
		for batchID := range sheetLayout.BatchCells {
			batchRegistry.AddBatch(types.SheetID(sheet), batchID)
		}
	}

	tableParser := parser.NewParser(file, layout)
	timetables, err := tableParser.Parse()
	if err != nil {
		return err
	}

	// Make & Populate TimetableRegistry
	tableRegistry := registry.NewTimetableRegistry()
	for _, timetable := range timetables {
		tableRegistry.AddTimetable(timetable.Batch, &timetable)
	}

	// Release workbook — no longer needed after parsing
	file.Close()
	runtime.GC()
	debug.FreeOSMemory()

	// API
	router := api.SetupRoutes(tableRegistry, batchRegistry, "frontend/dist")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	return nil
}
