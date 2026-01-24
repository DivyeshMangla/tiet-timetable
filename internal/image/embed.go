package image

import (
	"embed"
	"io"
)

//go:embed timetable-bg.png Inter-SemiBold.ttf
var FS embed.FS

const (
	Background = "timetable-bg.png"
	FontFile   = "Inter-SemiBold.ttf"
)

func GetBackground(name string) (io.ReadCloser, error) {
	file, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFont(name string) ([]byte, error) {
	return FS.ReadFile(name)
}
