package model

type ClassType int

const (
	LECTURE ClassType = iota
	TUTORIAL
	PRACTICAL
)

var classTypeSuffixes = map[ClassType]rune{
	LECTURE:   'L',
	TUTORIAL:  'T',
	PRACTICAL: 'P',
}

var suffixToClassType = map[rune]ClassType{
	'L': LECTURE,
	'T': TUTORIAL,
	'P': PRACTICAL,
}

func (ct ClassType) Suffix() rune {
	return classTypeSuffixes[ct]
}

func FromSuffix(suffix rune) *ClassType {
	if ct, ok := suffixToClassType[suffix]; ok {
		return &ct
	}

	return nil
}

func (ct ClassType) String() string {
	switch ct {
	case LECTURE:
		return "LECTURE"
	case TUTORIAL:
		return "TUTORIAL"
	case PRACTICAL:
		return "PRACTICAL"
	default:
		return "UNKNOWN"
	}
}
