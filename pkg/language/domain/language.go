package domain

import "strings"

type Language string

const (
	LanguageUnknown    Language = ""
	LanguageEnglish    Language = "english"
	LanguageVietnamese Language = "vietnamese"
)

func (s Language) String() string {
	switch s {
	case LanguageEnglish, LanguageVietnamese:
		return string(s)
	default:
		return ""
	}
}

func (s Language) IsValid() bool {
	return s != LanguageUnknown
}

func ToLanguage(value string) Language {
	switch strings.ToLower(value) {
	case LanguageEnglish.String():
		return LanguageEnglish
	case LanguageVietnamese.String():
		return LanguageVietnamese
	default:
		return LanguageUnknown
	}
}
