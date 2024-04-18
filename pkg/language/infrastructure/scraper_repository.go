package infrastructure

import (
	"github.com/namhq1989/bapbi-server/internal/scraper"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

type ScraperRepository struct {
	scraper *scraper.Scraper
}

func NewScraperRepository(scraper *scraper.Scraper) ScraperRepository {
	return ScraperRepository{
		scraper: scraper,
	}
}

func (r ScraperRepository) GetEnglishDictionaryData(ctx *appcontext.AppContext, term string) (*domain.EnglishDictionaryScraperData, error) {
	result, err := r.scraper.GetCambridgeEnglishDictionaryData(ctx, term)
	if err != nil {
		return nil, err
	}

	return &domain.EnglishDictionaryScraperData{
		Level:        result.Level,
		Phonetic:     result.Phonetic,
		PartOfSpeech: result.PartOfSpeech,
		AudioURL:     result.AudioURL,
	}, nil
}
