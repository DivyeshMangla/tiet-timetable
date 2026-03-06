package parser

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/xuri/excelize/v2"
)

type Worksheet struct {
	file       *excelize.File
	sheet      string
	mergeIndex map[[2]int]*utils.MergedRegion
}

func NewWorksheet(file *excelize.File, sheet string) (*Worksheet, error) {
	index, err := buildMergeIndex(file, sheet)
	if err != nil {
		return nil, err
	}

	return &Worksheet{
		file:       file,
		sheet:      sheet,
		mergeIndex: index,
	}, nil
}

func (ws *Worksheet) Cell(row, col int) (string, error) {
	cell, err := excelize.CoordinatesToCellName(col+1, row+1)
	if err != nil {
		return "", err
	}
	return ws.file.GetCellValue(ws.sheet, cell)
}

func (ws *Worksheet) HorizontalMergedRegion(row, col int) (*utils.MergedRegion, bool) {
	r, ok := ws.mergeIndex[[2]int{row, col}]
	if !ok {
		return nil, false
	}

	if r.StartRow == r.EndRow && r.StartCol != r.EndCol {
		return r, true
	}

	return nil, false
}

func buildMergeIndex(file *excelize.File, sheet string) (map[[2]int]*utils.MergedRegion, error) {
	mergeCells, err := file.GetMergeCells(sheet)
	if err != nil {
		return nil, err
	}

	index := map[[2]int]*utils.MergedRegion{}

	for _, mc := range mergeCells {
		sc, sr, _ := excelize.CellNameToCoordinates(mc.GetStartAxis())
		ec, er, _ := excelize.CellNameToCoordinates(mc.GetEndAxis())

		r := &utils.MergedRegion{
			StartRow: sr - 1,
			EndRow:   er - 1,
			StartCol: sc - 1,
			EndCol:   ec - 1,
		}

		for row := r.StartRow; row <= r.EndRow; row++ {
			for col := r.StartCol; col <= r.EndCol; col++ {
				index[[2]int{row, col}] = r
			}
		}
	}

	return index, nil
}
