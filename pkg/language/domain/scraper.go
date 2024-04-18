package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type ScraperRepository interface {
	GetEnglishDictionaryData(ctx *appcontext.AppContext, term string) (*EnglishDictionaryScraperData, error)
}

type EnglishDictionaryScraperData struct {
	Level        string
	Phonetic     string
	PartOfSpeech string
	AudioURL     string
}
