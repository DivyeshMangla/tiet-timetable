package utils

import "strings"

func IsValid(value string) bool {
	return strings.TrimSpace(value) != ""
}

func FirstLine(value string) string {
	if idx := strings.IndexAny(value, "\r\n"); idx != -1 {
		return value[:idx]
	}
	return value
}
