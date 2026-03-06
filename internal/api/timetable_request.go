package api

import (
	"fmt"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type Subjects struct {
	Code  string `json:"code"`
	Alias string `json:"alias,omitempty"`
}

type FormattedTimetableRequest struct {
	Batch    string     `json:"batch"`
	Subjects []Subjects `json:"subjects"`
}

func BuildRenderableTimetable(timetable *types.Timetable, subjects []Subjects) *types.RenderableTimetable {
	subjectSet := make(map[types.SubjectCode]struct{}, len(subjects))
	aliasMap := make(map[types.SubjectCode]string, len(subjects))

	for _, s := range subjects {
		code := types.SubjectCode(s.Code)
		subjectSet[code] = struct{}{}
		aliasMap[code] = s.Alias
	}

	result := &types.RenderableTimetable{
		Batch: timetable.Batch,
		Days:  make(map[types.Day][]types.RenderInfo),
	}

	for day, slots := range timetable.Days {
		for _, slot := range slots {
			for _, class := range slot.Classes {

				if _, ok := subjectSet[class.SubjectCode]; !ok {
					continue
				}

				label := string(class.SubjectCode)
				if alias := aliasMap[class.SubjectCode]; alias != "" {
					label = alias
				}

				text := fmt.Sprintf("%s - %s", label, class.Room)

				result.Days[day] = append(result.Days[day], types.RenderInfo{
					Start:     slot.Start,
					End:       slot.End,
					ClassType: class.SubjectCode.ClassType(),
					Text:      text,
				})
			}
		}
	}

	return result
}
