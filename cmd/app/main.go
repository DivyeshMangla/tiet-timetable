package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DivyeshMangla/tiet-timetable/internal/api"
	"github.com/DivyeshMangla/tiet-timetable/internal/config"
	"github.com/DivyeshMangla/tiet-timetable/internal/io"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/xuri/excelize/v2"
)

const configFile = "config.yml"
const serverPort = ":8080"

func main() {
	loader := config.NewConfigLoader(configFile)

	handled, err := loader.HandleInitFlag(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to handle init flag: %v", err)
	}
	if handled {
		return
	}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	workbook, err := loadFromURL(cfg.TimetableURL)
	if err != nil {
		log.Fatalf("Failed to load workbook: %v", err)
	}
	defer workbook.Close()

	p, err := parser.NewParser(workbook)
	if err != nil {
		log.Fatalf("Failed to create parser: %v", err)
	}

	fmt.Println("Populating registry...")
	reg := registry.NewTimetableRegistry()
	if err := registry.PopulateFromParser(reg, p); err != nil {
		log.Fatalf("Failed to populate registry: %v", err)
	}

	router := api.SetupRoutes(reg)

	fmt.Printf("Server starting on http://localhost%s\n", serverPort)
	if err := http.ListenAndServe(serverPort, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loadFromURL(url string) (*excelize.File, error) {
	reader, err := io.Download(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer reader.Close()

	workbook, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open workbook: %w", err)
	}
	return workbook, nil
}
