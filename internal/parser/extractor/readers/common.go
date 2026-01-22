package readers

import "strings"

func isValid(value string) bool {
	return strings.TrimSpace(value) != ""
}
