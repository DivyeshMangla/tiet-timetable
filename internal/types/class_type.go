package types

type ClassType int

const (
	LECTURE ClassType = iota
	TUTORIAL
	PRACTICAL
)

func (s SubjectCode) ClassType() ClassType {
	if len(s) == 0 {
		return LECTURE
	}

	switch s[len(s)-1] {
	case 'L':
		return LECTURE
	case 'T':
		return TUTORIAL
	case 'P':
		return PRACTICAL
	default:
		return LECTURE
	}
}
