package model

import "github.com/DivyeshMangla/tiet-timetable/internal/types"

type Subject struct {
	Code types.SubjectCode
	Name types.SubjectName
	Abbr types.SubjectAbbr
}

func NewSubject(code types.SubjectCode, name types.SubjectName, abbr types.SubjectAbbr) Subject {
	if abbr == "" {
		abbr = types.SubjectAbbr(name)
	}
	return Subject{
		Code: code,
		Name: name,
		Abbr: abbr,
	}
}
