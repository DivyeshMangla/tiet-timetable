package image

import (
	"embed"
	"io"
)

//go:embed timetable-bg.png
var FS embed.FS

const (
	Background = "timetable-bg.png"
)

func GetBackground(name string) (io.ReadCloser, error) {
	file, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}
