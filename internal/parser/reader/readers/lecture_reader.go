package readers

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/excel"
	"github.com/DivyeshMangla/tiet-timetable/internal/parser/utils"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
	"regexp"
	"strings"
)

const minHorizontalMergeWidth = 3

var teacherPattern = regexp.MustCompile(`^[A-Za-z0-9. -]+$`)
var subjectCodePattern = regexp.MustCompile(`^([A-Z]{3}\d{3}|[A-Z]{5}\d)[LTP]$`)

func parseSubjectCode(code string) (types.SubjectCode, types.ClassType) {
	suffix := code[len(code)-1]
	stripped := code[:len(code)-1]

	var ct types.ClassType
	switch suffix {
	case 'L':
		ct = types.LECTURE
	case 'T':
		ct = types.TUTORIAL
	case 'P':
		ct = types.PRACTICAL
	}

	return types.SubjectCode(stripped), ct
}

type LectureReader struct{}

func (l LectureReader) Read(ws *excel.Worksheet, start types.TimeSlot, row, col int) (*types.ClassSlot, bool) {
	region, found := ws.HorizontalMergedRegion(row, col)
	if !found || !isWideEnough(region) {
		return nil, false
	}

	subjectCodeValue, err := ws.Cell(row, region.StartCol)
	if err != nil {
		return nil, false
	}

	subjectMatcher := utils.NewValueMatcher(CleanCell(subjectCodeValue), subjectCodePattern)
	if !subjectMatcher.Valid() {
		return nil, false
	}

	roomValue, err := ws.Cell(row+1, region.StartCol)
	if err != nil {
		return nil, false
	}
	roomMatcher := utils.NewValueMatcher(CleanCell(roomValue), nil)

	teacherValue, err := ws.Cell(row+1, region.EndCol)
	if err != nil {
		return nil, false
	}

	teacherMatcher := utils.NewValueMatcher(CleanCell(teacherValue), teacherPattern)
	if !teacherMatcher.Valid() {
		return nil, false
	}

	subjects := subjectMatcher.Values()
	rooms := roomMatcher.Values()
	teachers := teacherMatcher.Values()

	n := len(subjects)

	if len(rooms) < n || len(teachers) < n {
		return nil, false
	}

	classes := make([]types.Class, n)

	for i := 0; i < n; i++ {
		code, ct := parseSubjectCode(subjects[i])
		classes[i] = types.Class{
			SubjectCode: code,
			ClassType:   ct,
			Room:        types.Room(rooms[i]),
			Teacher:     types.Teacher(teachers[i]),
		}
	}

	return &types.ClassSlot{
		Start:   start,
		End:     start,
		Classes: classes,
	}, true
}

func isWideEnough(region *utils.MergedRegion) bool {
	return (region.EndCol - region.StartCol) >= minHorizontalMergeWidth
}

func CleanCell(value string) string {
	if idx := strings.IndexAny(value, "\r\n"); idx != -1 {
		value = value[:idx]
	}
	return strings.TrimSpace(value)
}
