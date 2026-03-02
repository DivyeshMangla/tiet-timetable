package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/api"
	"github.com/DivyeshMangla/tiet-timetable/internal/io"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/registry"
	"github.com/xuri/excelize/v2"
)

const defaultTimetableURL = "https://www.thapar.edu/upload/files/UG,%20PG%20TIME%20TABLE%20JAN%20TO%20MAY%202026.xlsx"

func main() {
	timetableURL := os.Getenv("TIMETABLE_URL")
	if timetableURL == "" {
		timetableURL = defaultTimetableURL
	}

	workbook, err := loadFromURL(timetableURL)
	if err != nil {
		log.Fatalf("Failed to load workbook: %v", err)
	}
	if err := TrimSheetNames(workbook); err != nil {
		log.Fatalf("Failed to trim sheet names: %v", err)
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

	router := api.SetupRoutes(reg, "frontend/dist")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
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

func TrimSheetNames(file *excelize.File) error {
	sheetList := file.GetSheetList()

	for _, oldName := range sheetList {
		visible, err := file.GetSheetVisible(oldName)
		if err != nil {
			return err
		}
		if !visible {
			continue
		}

		newName := strings.TrimSpace(oldName)
		if newName != oldName {
			if err := file.SetSheetName(oldName, newName); err != nil {
				return err
			}
		}
	}

	return nil
}
