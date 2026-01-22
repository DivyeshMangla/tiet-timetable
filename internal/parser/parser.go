package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/extractor"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type Parser struct {
	cache          *Cache
	classExtractor *extractor.ClassExtractor
	file           *excelize.File
}

func NewParser(file *excelize.File) (*Parser, error) {
	start := time.Now()

	cache, err := FromWorkbook(file)
	if err != nil {
		return nil, fmt.Errorf("failed to build cache: %w", err)
	}

	totalBatches := 0
	for _, batches := range cache.batches {
		totalBatches += len(batches)
	}

	fmt.Printf("Parsed %d sheets with %d batches in %.2fms\n",
		len(cache.batches), totalBatches, float64(time.Since(start).Nanoseconds())/1e6)

	return &Parser{
		cache:          cache,
		classExtractor: extractor.NewClassExtractor(),
		file:           file,
	}, nil
}

func (p *Parser) GetSheetByName(sheetName string) (types.SheetID, bool) {
	normalized := normalizeWhitespace(sheetName)
	for sheetID := range p.cache.batches {
		if normalizeWhitespace(string(sheetID)) == normalized {
			return sheetID, true
		}
	}
	return "", false
}

func (p *Parser) SheetNames() []string {
	names := make([]string, 0, len(p.cache.batches))
	for sheetID := range p.cache.batches {
		names = append(names, string(sheetID))
	}
	return names
}

func (p *Parser) BatchNames(sheetName string) []types.BatchID {
	sheetID, ok := p.GetSheetByName(sheetName)
	if !ok {
		return nil
	}

	batches := p.cache.batches[sheetID]
	if len(batches) == 0 {
		return nil
	}

	ids := make([]types.BatchID, 0, len(batches))
	for id := range batches {
		ids = append(ids, id)
	}
	return ids
}

func (p *Parser) GetTimetable(sheetName, batchName string) (model.Timetable, error) {
	sheetID, ok := p.GetSheetByName(sheetName)
	if !ok {
		return model.Timetable{}, fmt.Errorf("sheet not found: %s", sheetName)
	}

	batchID := types.BatchID(batchName)
	batchLocation, ok := p.cache.batches[sheetID][batchID]
	if !ok {
		return model.Timetable{}, fmt.Errorf("batch not found: %s", batchName)
	}

	daySlotsList, ok := p.cache.daySlots[sheetID]
	if !ok {
		return model.Timetable{}, fmt.Errorf("day slots not found for sheet: %s", sheetName)
	}

	var entries []model.TimetableEntry
	for _, daySlots := range daySlotsList {
		p.processDay(string(sheetID), batchLocation.Col, daySlots, &entries)
	}

	return model.NewTimetable(entries), nil
}

func (p *Parser) processDay(sheetName string, batchCol int, daySlots DaySlots, entries *[]model.TimetableEntry) {
	for timeSlot, loc := range daySlots.Slots {
		p.processTimeSlot(sheetName, batchCol, daySlots, timeSlot, loc.Row, entries)
	}
}

func (p *Parser) processTimeSlot(
	sheetName string,
	batchCol int,
	daySlots DaySlots,
	timeSlot model.TimeSlot,
	row int,
	entries *[]model.TimetableEntry,
) {
	classInfo := p.classExtractor.Extract(p.file, sheetName, row, batchCol)
	if classInfo == nil {
		return
	}

	region, found := utils.GetVerticalMergedRegion(p.file, sheetName, row, batchCol)
	isBlock := found && region.StartRow != region.EndRow

	if isBlock {
		p.processBlockClass(region, row, daySlots, *classInfo, entries)
		return
	}

	*entries = append(*entries, model.TimetableEntry{
		Day:       daySlots.Day,
		TimeSlot:  timeSlot,
		ClassInfo: *classInfo,
	})
}

func (p *Parser) processBlockClass(
	region utils.MergedRegion,
	currentRow int,
	daySlots DaySlots,
	classInfo model.ClassInfo,
	entries *[]model.TimetableEntry,
) {
	if currentRow != region.StartRow {
		return
	}

	for slot, loc := range daySlots.Slots {
		if loc.Row >= region.StartRow && loc.Row <= region.EndRow {
			*entries = append(*entries, model.TimetableEntry{
				Day:       daySlots.Day,
				TimeSlot:  slot,
				ClassInfo: classInfo,
			})
		}
	}
}

func normalizeWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r != ' ' && r != '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
