package registry

import (
	"fmt"
	"time"

	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
)

func PopulateFromParser(reg *TimetableRegistry, p *parser.Parser) error {
	sheetNames := p.SheetNames()
	totalBatches := 0

	for _, sheetName := range sheetNames {
		sheetID, ok := p.GetSheetByName(sheetName)
		if !ok {
			continue
		}

		batchNames := p.BatchNames(sheetName)
		if len(batchNames) == 0 {
			continue
		}

		for _, batchID := range batchNames {
			start := time.Now()
			timetable, err := p.GetTimetable(sheetName, string(batchID))
			if err != nil {
				return fmt.Errorf("failed to get timetable for sheet %q batch %q: %w", sheetName, batchID, err)
			}

			elapsed := time.Since(start)
			entryCount := len(timetable.Entries)

			fmt.Printf("Loaded timetable for batch %s (%d entries) in %.2fms\n", batchID, entryCount, float64(elapsed.Nanoseconds())/1e6)
			reg.RegisterBatch(sheetID, batchID, timetable)
			totalBatches++
		}
	}

	fmt.Printf("\nSuccessfully loaded %d batches from %d sheets\n", totalBatches, len(sheetNames))
	return nil
}
