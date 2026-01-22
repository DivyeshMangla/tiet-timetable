package registry

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

func PopulateFromParser(reg *TimetableRegistry, p *parser.Parser) {
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
			reg.RegisterBatchMetadata(sheetID, batchID)
		}
	}

	reg.parser = p
}
