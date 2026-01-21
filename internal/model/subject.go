package model

type Subject struct {
	Code string
	Name string
	Abbr string
}

func NewSubject(code, name, abbr string) Subject {
	if abbr == "" {
		abbr = name // Default to name if abbreviation is empty
	}

	return Subject{
		Code: code,
		Name: name,
		Abbr: abbr,
	}
}
