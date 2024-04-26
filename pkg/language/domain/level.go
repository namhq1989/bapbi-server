package domain

import (
	"strings"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
)

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

func GetMinWordsBasedOnLevel(level Level) int {
	index := manipulation.RandomIntInRange(0, 5)

	switch level {
	case LevelBeginner:
		return []int{30, 35, 40, 45, 50}[index]
	case LevelIntermediate:
		return []int{80, 90, 100, 110, 120}[index]
	case LevelAdvanced:
		return []int{150, 170, 190, 210, 230}[index]
	default:
		return 50
	}
}
