package utils

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

type SheetCache struct {
	regions []MergedRegion
	byCell  map[string]*MergedRegion
}

type MergeCache struct {
	mu     sync.RWMutex
	sheets map[string]*SheetCache
}

var (
	instance *MergeCache
	once     sync.Once
)

func GetMergeCacheInstance() *MergeCache {
	once.Do(func() {
		instance = &MergeCache{
			sheets: make(map[string]*SheetCache),
		}
	})
	return instance
}

func (mc *MergeCache) getOrBuildCache(file *excelize.File, sheetName string) (*SheetCache, error) {
	mc.mu.RLock()
	cache, exists := mc.sheets[sheetName]
	mc.mu.RUnlock()

	if exists {
		return cache, nil
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	if cache, exists := mc.sheets[sheetName]; exists {
		return cache, nil
	}

	regions, err := parseMergedRegions(file, sheetName)
	if err != nil {
		return nil, err
	}

	cache = &SheetCache{
		regions: regions,
		byCell:  make(map[string]*MergedRegion),
	}

	for i := range regions {
		r := &regions[i]
		for row := r.StartRow; row <= r.EndRow; row++ {
			for col := r.StartCol; col <= r.EndCol; col++ {
				key := fmt.Sprintf("%d,%d", row, col)
				cache.byCell[key] = r
			}
		}
	}

	mc.sheets[sheetName] = cache
	return cache, nil
}

func (mc *MergeCache) ClearSheet(sheetName string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	delete(mc.sheets, sheetName)
}

func (mc *MergeCache) ClearAll() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.sheets = make(map[string]*SheetCache)
}
