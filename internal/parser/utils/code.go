package utils

import (
	"github.com/DivyeshMangla/tiet-timetable/internal/model"
	"github.com/DivyeshMangla/tiet-timetable/internal/types"
)

type ParsedCode struct {
	Code      types.SubjectCode
	ClassType model.ClassType
}

func ParseCode(code string) (ParsedCode, bool) {
	if code == "" {
		return ParsedCode{}, false
	}

	if len(code) == 0 {
		return ParsedCode{}, false
	}

	lastChar := rune(code[len(code)-1])
	codeWithoutSuffix := code[:len(code)-1]

	classType := model.FromSuffix(lastChar)
	if classType == nil {
		return ParsedCode{}, false
	}

	return ParsedCode{
		Code:      types.SubjectCode(codeWithoutSuffix),
		ClassType: *classType,
	}, true
}
