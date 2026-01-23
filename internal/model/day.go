package model

type Day int

const (
	MON Day = iota
	TUE
	WED
	THU
	FRI
)

func (d Day) String() string {
	switch d {
	case MON:
		return "MON"
	case TUE:
		return "TUE"
	case WED:
		return "WED"
	case THU:
		return "THU"
	case FRI:
		return "FRI"
	default:
		return "UNKNOWN"
	}
}
