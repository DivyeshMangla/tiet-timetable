package image

const (
	TextColorHex       = "000000"
	BackgroundColorHex = "191816"
	CellColorHex       = "FFFBF0"
	LectureColorHex    = "F7D89C"
	TutorialColorHex   = "FF938D"
	PracticalColorHex  = "BBE4F7"
)

var (
	TextColor       = HexToRGBA(TextColorHex)
	BackgroundColor = HexToRGBA(BackgroundColorHex)
	CellColor       = HexToRGBA(CellColorHex)
	LectureColor    = HexToRGBA(LectureColorHex)
	TutorialColor   = HexToRGBA(TutorialColorHex)
	PracticalColor  = HexToRGBA(PracticalColorHex)
)
