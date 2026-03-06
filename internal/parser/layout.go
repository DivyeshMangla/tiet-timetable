package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type DayLayout struct {
	Day           types.Day
	TimeSlotCells map[types.TimeSlot]types.CellLocation
}
type SheetLayout struct {
	Days       []DayLayout
	BatchCells map[types.BatchID]types.CellLocation
}

type WorkbookLayout struct {
	Sheets map[string]*SheetLayout
}

type WorkbookLayoutBuilder struct {
	file *excelize.File
}

func NewWorkbookLayoutBuilder(file *excelize.File) *WorkbookLayoutBuilder {
	return &WorkbookLayoutBuilder{file: file}
}

func (b *WorkbookLayoutBuilder) Build() (*WorkbookLayout, error) {
	sheets := b.file.GetSheetList()

	result := &WorkbookLayout{
		Sheets: make(map[string]*SheetLayout),
	}

	for _, sheet := range sheets {
		layout, err := buildSheetLayout(b.file, sheet)
		if err != nil {
			return nil, err
		}
		if layout != nil {
			result.Sheets[sheet] = layout
		}
	}

	return result, nil
}

var batchRegex = regexp.MustCompile(`^\d[A-Z]\d[A-Z]$`)

const maxRowsToScan = 300

func findDayRow(rows [][]string) int {
	for r, row := range rows {
		if len(row) > 0 && strings.EqualFold(strings.TrimSpace(row[0]), "day") {
			return r
		}
	}
	return -1
}

func buildDayLayouts(file *excelize.File, sheet string) ([]DayLayout, error) {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	firstRow := findDayRow(rows)
	if firstRow == -1 {
		return nil, nil
	}
	firstCol := 1

	result := make([]DayLayout, 0, len(types.Weekdays))
	slots := make(map[types.TimeSlot]types.CellLocation)

	dayIndex := -1
	lastSlot := -1

	// Scan rows for slot numbers; slot "1" marks the start of a new day
	for row := firstRow; row < firstRow+maxRowsToScan; row++ {
		value, _ := utils.GetCell(file, sheet, row, firstCol)

		num, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || num == lastSlot {
			continue
		}
		lastSlot = num

		if num == 1 {
			// Commit previous day before starting the next
			if dayIndex >= 0 {
				result = append(result, DayLayout{Day: types.Weekdays[dayIndex], TimeSlotCells: slots})
			}

			dayIndex++
			if dayIndex >= len(types.Weekdays) {
				break
			}

			slots = make(map[types.TimeSlot]types.CellLocation)
		}

		if slot, err := types.TimeSlotFromNumber(num); err == nil {
			slots[slot] = types.CellLocation{Row: row, Col: firstCol}
		}
	}

	// Append the last day's slots
	if dayIndex >= 0 && dayIndex < len(types.Weekdays) {
		result = append(result, DayLayout{Day: types.Weekdays[dayIndex], TimeSlotCells: slots})
	}

	return result, nil
}

func buildSheetLayout(file *excelize.File, sheet string) (*SheetLayout, error) {
	days, err := buildDayLayouts(file, sheet)
	if err != nil || len(days) == 0 {
		return nil, err
	}

	batches, err := findBatchCells(file, sheet)
	if err != nil {
		return nil, err
	}

	return &SheetLayout{
		Days:       days,
		BatchCells: batches,
	}, nil
}

func findBatchCells(file *excelize.File, sheet string) (map[types.BatchID]types.CellLocation, error) {
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	dayRow := findDayRow(rows)
	if dayRow == -1 || dayRow >= len(rows) {
		return nil, nil
	}

	result := make(map[types.BatchID]types.CellLocation)

	// Scan across the "day" row looking for batch IDs
	for col, value := range rows[dayRow] {
		cell := strings.TrimSpace(value)
		if !batchRegex.MatchString(cell) {
			continue
		}

		id := types.BatchID(normalizeBatchName(cell))
		result[id] = types.CellLocation{Row: dayRow, Col: col}
	}

	return result, nil
}

func normalizeBatchName(raw string) string {
	if len(raw) != 4 {
		return raw
	}
	n := int(raw[3]-'A') + 1
	return raw[:3] + strconv.Itoa(n)
}
