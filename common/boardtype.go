package common

type BoardType byte

func NewBoardTypeFromString(s string) BoardType {
	if s == "dcb" {
		return DCBBoard
	}
	if s == "gov" {
		return GOVBoard
	}
	return 255
}

func (boardType *BoardType) Bytes() []byte {
	x := byte(*boardType)
	return []byte{x}
}
