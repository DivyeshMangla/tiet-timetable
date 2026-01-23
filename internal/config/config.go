package config

import "github.com/DivyeshMangla/tiet-timetable/internal/model"

type Config struct {
	Subjects     []model.Subject
	TimetableURL string
}
