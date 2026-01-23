package model

import "github.com/DivyeshMangla/tiet-timetable/internal/types"

type ClassInfo struct {
	SubjectCode types.SubjectCode
	ClassType   ClassType
	Room        types.Room
	Teacher     types.Teacher
	IsBlock     bool
}
