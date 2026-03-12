package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/reader/readers"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"github.com/xuri/excelize/v2"
)

type Parser struct {
	workbook *excelize.File
	layout   *WorkbookLayout
}

func NewParser(file *excelize.File, layout *WorkbookLayout) *Parser {
	return &Parser{
		workbook: file,
		layout:   layout,
	}
}

func (p *Parser) Parse() ([]types.Timetable, error) {
	var result []types.Timetable

	extractor := readers.NewClassExtractor()

	for sheetName, sheetLayout := range p.layout.Sheets {

		ws, err := excel.NewWorksheet(p.workbook, sheetName)
		if err != nil {
			return nil, err
		}

		for batchID, batchCell := range sheetLayout.BatchCells {

			timetable := types.Timetable{
				Batch: batchID,
				Days:  make(map[types.Day][]types.ClassSlot),
			}

			for _, dayLayout := range sheetLayout.Days {

				for slot, slotCell := range dayLayout.TimeSlotCells {

					row := slotCell.Row
					col := batchCell.Col

					classSlot := extractor.Extract(ws, types.TimeSlot(slot), row, col)
					if classSlot == nil {
						continue
					}

					timetable.Days[dayLayout.Day] =
						append(timetable.Days[dayLayout.Day], *classSlot)
				}
			}

			result = append(result, timetable)
		}
	}

	return result, nil
}
