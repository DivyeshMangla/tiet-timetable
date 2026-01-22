package registry

import (
	"fmt"

	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
)

func PopulateFromParser(reg *TimetableRegistry, p *parser.Parser) error {
	sheetNames := p.SheetNames()

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
			timetable, err := p.GetTimetable(sheetName, string(batchID))
			if err != nil {
				return fmt.Errorf("failed to get timetable for sheet %q batch %q: %w", sheetName, batchID, err)
			}

			reg.RegisterBatch(sheetID, batchID, timetable)
		}
	}

	return nil
}
