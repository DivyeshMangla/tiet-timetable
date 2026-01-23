package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

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
	cache := GetMergeCacheInstance()
	sheetCache, err := cache.getOrBuildCache(file, sheetName)
	if err != nil {
		return nil, err
	}
	return sheetCache.regions, nil
}

func GetHorizontalMergedRegion(file *excelize.File, sheetName string, row, col int) (MergedRegion, bool) {
	cache := GetMergeCacheInstance()
	sheetCache, err := cache.getOrBuildCache(file, sheetName)
	if err != nil {
		return MergedRegion{}, false
	}

	key := fmt.Sprintf("%d,%d", row, col)
	if region, ok := sheetCache.byCell[key]; ok && region.IsHorizontal() {
		return *region, true
	}

	return MergedRegion{}, false
}

func parseMergedRegions(file *excelize.File, sheetName string) ([]MergedRegion, error) {
	mergeCells, err := file.GetMergeCells(sheetName)
	if err != nil {
		return nil, err
	}

	regions := make([]MergedRegion, 0, len(mergeCells))

	for _, mc := range mergeCells {
		startCol, startRow, err := excelize.CellNameToCoordinates(mc.GetStartAxis())
		if err != nil {
			continue
		}

		endCol, endRow, err := excelize.CellNameToCoordinates(mc.GetEndAxis())
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
