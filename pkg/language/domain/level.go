package domain

import "strings"

type Level string

const (
	LevelUnknown      Level = ""
	LevelBeginner     Level = "beginner"
	LevelIntermediate Level = "intermediate"
	LevelAdvanced     Level = "advanced"
)

var ListLevels = []Level{
	LevelBeginner,
	LevelIntermediate,
	LevelAdvanced,
}

func (s Level) String() string {
	switch s {
	case LevelBeginner, LevelIntermediate, LevelAdvanced:
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
	case LevelBeginner.String():
		return LevelBeginner
	case LevelIntermediate.String():
		return LevelIntermediate
	case LevelAdvanced.String():
		return LevelAdvanced
	default:
		return LevelUnknown
	}
}
