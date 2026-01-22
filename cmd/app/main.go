package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DivyeshMangla/tiet-timetable/internal/config"
	"github.com/DivyeshMangla/tiet-timetable/internal/io"
	"github.com/xuri/excelize/v2"
)

const configFile = "config.yml"

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

	_ = workbook
	fmt.Println("Timetable loaded successfully")
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
