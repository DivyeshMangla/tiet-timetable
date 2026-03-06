package utils

import "regexp"

type ValueMatcher struct {
	values []string
	regex  *regexp.Regexp
}

func NewValueMatcher(input string, pattern *regexp.Regexp) *ValueMatcher {
	parts := splitValues(input)

	vm := &ValueMatcher{
		values: parts,
		regex:  pattern,
	}

	return vm
}

func splitValues(input string) []string {
	if input == "" {
		return nil
	}

	var result []string
	start := 0

	for i := 0; i < len(input); i++ {
		if input[i] == '/' {
			result = append(result, input[start:i])
			start = i + 1
		}
	}

	result = append(result, input[start:])
	return result
}

func (v *ValueMatcher) Valid() bool {
	if v.regex == nil {
		return true
	}

	for _, val := range v.values {
		if !v.regex.MatchString(val) {
			return false
		}
	}

	return true
}

func (v *ValueMatcher) Values() []string {
	return v.values
}

func (v *ValueMatcher) HasOneValue() bool {
	return len(v.values) == 1
}

func (v *ValueMatcher) First() string {
	if len(v.values) == 0 {
		return ""
	}
	return v.values[0]
}
