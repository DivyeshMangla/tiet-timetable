package extractor

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

var batchRegex = regexp.MustCompile(`^\d[A-Z]\d[A-Z]$`)

type BatchExtractor struct {
	file      *excelize.File
	sheetName string
}

func NewBatchExtractor(file *excelize.File, sheetName string) *BatchExtractor {
	return &BatchExtractor{
		file:      file,
		sheetName: sheetName,
	}
}

func (be *BatchExtractor) Extract() (map[types.BatchID]parser.CellLocation, error) {
	dayRow, _, found := parser.FindCellInFirstColumn(be.file, be.sheetName, "day")
	if !found {
		return nil, fmt.Errorf("could not find 'day' cell in the first column")
	}

	rows, err := be.file.GetRows(be.sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if dayRow >= len(rows) {
		return nil, fmt.Errorf("day row out of bounds")
	}

	batches := make(map[types.BatchID]parser.CellLocation)
	dayRowData := rows[dayRow]

	for col := 0; col < len(dayRowData); col++ {
		cellValue := strings.TrimSpace(dayRowData[col])
		if cellValue == "" {
			continue
		}

		if batchRegex.MatchString(cellValue) {
			normalized := normalizeBatchName(cellValue)
			batches[types.BatchID(normalized)] = parser.CellLocation{
				Row: dayRow,
				Col: col,
			}
		}
	}

	return sortBatches(batches), nil
}

func extractNumber(batchName string) int {
	if len(batchName) < 3 {
		return 0
	}
	num, _ := strconv.Atoi(batchName[2:3])
	return num
}

func normalizeBatchName(rawName string) string {
	if len(rawName) != 4 {
		return rawName
	}

	lastLetter := rawName[3]
	position := int(lastLetter-'A') + 1

	return rawName[:3] + strconv.Itoa(position)
}

func sortBatches(batches map[types.BatchID]parser.CellLocation) map[types.BatchID]parser.CellLocation {
	type batchEntry struct {
		id       types.BatchID
		location parser.CellLocation
	}

	entries := make([]batchEntry, 0, len(batches))
	for id, loc := range batches {
		entries = append(entries, batchEntry{id: id, location: loc})
	}

	sort.Slice(entries, func(i, j int) bool {
		num1 := extractNumber(string(entries[i].id))
		num2 := extractNumber(string(entries[j].id))
		if num1 != num2 {
			return num1 < num2
		}
		return string(entries[i].id) < string(entries[j].id)
	})

	sorted := make(map[types.BatchID]parser.CellLocation, len(entries))
	for _, entry := range entries {
		sorted[entry.id] = entry.location
	}

	return sorted
}
