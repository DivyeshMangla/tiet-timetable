package main

import (
	"fmt"
	"github.com/DivyeshMangla/tiet-timetable/internal"
	"log"
	"os"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/io"
	"github.com/xuri/excelize/v2"
)

const defaultTimetableURL = "https://www.thapar.edu/upload/files/UG,%20PG%20TIME%20TABLE%20JAN%20TO%20MAY%202026.xlsx"

func main() {
	println("Starting timetable server...")
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

	if err := internal.Bootstrap(workbook); err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
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
