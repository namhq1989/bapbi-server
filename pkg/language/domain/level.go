package domain

import "strings"

type Level string

const (
	LevelUnknown Level = ""
	LevelA       Level = "A"
	LevelB       Level = "B"
	LevelC       Level = "C"
)

func (s Level) String() string {
	switch s {
	case LevelA, LevelB, LevelC:
		return string(s)
	default:
		return ""
	}
}

func (s Level) IsValid() bool {
	return s != LevelUnknown
}

func ToLevel(value string) Level {
	switch strings.ToLower(value) {
	case LevelA.String():
		return LevelA
	case LevelC.String():
		return LevelB
	case LevelC.String():
		return LevelC
	default:
		return LevelUnknown
	}
}
