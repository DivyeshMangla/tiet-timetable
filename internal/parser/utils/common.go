package utils

import "strings"

func IsValid(value string) bool {
	return strings.TrimSpace(value) != ""
}
