package utils

import "github.com/xuri/excelize/v2"

type MergedRegion struct {
	StartRow int
	EndRow   int
	StartCol int
	EndCol   int
}

func (m MergedRegion) IsInRange(row, col int) bool {
	return row >= m.StartRow && row <= m.EndRow &&
		col >= m.StartCol && col <= m.EndCol
}

func (m MergedRegion) IsVertical() bool {
	return m.StartRow != m.EndRow && m.StartCol == m.EndCol
}

func (m MergedRegion) IsHorizontal() bool {
	return m.StartRow == m.EndRow && m.StartCol != m.EndCol
}

func GetMergedRegions(file *excelize.File, sheetName string) ([]MergedRegion, error) {
	mergeCells, err := file.GetMergeCells(sheetName)
	if err != nil {
		return nil, err
	}

	regions := make([]MergedRegion, 0, len(mergeCells))

	for _, mc := range mergeCells {
		startRow, startCol, err := excelize.CellNameToCoordinates(mc.GetStartAxis())
		if err != nil {
			continue
		}

		endRow, endCol, err := excelize.CellNameToCoordinates(mc.GetEndAxis())
		if err != nil {
			continue
		}

		regions = append(regions, MergedRegion{
			StartRow: startRow - 1,
			EndRow:   endRow - 1,
			StartCol: startCol - 1,
			EndCol:   endCol - 1,
		})
	}

	return regions, nil
}

func GetVerticalMergedRegion(file *excelize.File, sheetName string, row, col int) (MergedRegion, bool) {
	regions, err := GetMergedRegions(file, sheetName)
	if err != nil {
		return MergedRegion{}, false
	}

	for _, region := range regions {
		if region.IsVertical() && region.IsInRange(row, col) {
			return region, true
		}
	}

	return MergedRegion{}, false
}

func GetHorizontalMergedRegion(file *excelize.File, sheetName string, row, col int) (MergedRegion, bool) {

	mergeCells, err := file.GetMergeCells(sheetName)
	if err != nil {
		return MergedRegion{}, false
	}

	r := row + 1
	c := col + 1

	for _, mc := range mergeCells {
		startCol, startRow, _ := excelize.CellNameToCoordinates(mc.GetStartAxis())
		endCol, endRow, _ := excelize.CellNameToCoordinates(mc.GetEndAxis())

		// horizontal merge = spans columns, not rows
		if startRow == endRow &&
			r == startRow &&
			c >= startCol && c <= endCol {

			return MergedRegion{
				StartRow: startRow - 1,
				EndRow:   endRow - 1,
				StartCol: startCol - 1,
				EndCol:   endCol - 1,
			}, true
		}
	}

	return MergedRegion{}, false
}
