package query

import (
	"sync"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SearchTermHandler struct {
	termRepository    domain.TermRepository
	openaiRepository  domain.OpenAIRepository
	scraperRepository domain.ScraperRepository
}

func NewSearchTermHandler(
	termRepository domain.TermRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
) SearchTermHandler {
	return SearchTermHandler{
		termRepository:    termRepository,
		openaiRepository:  openaiRepository,
		scraperRepository: scraperRepository,
	}
}

func (h SearchTermHandler) SearchTerm(ctx *appcontext.AppContext, performerID string, req dto.SearchTermRequest) (*dto.SearchTermResponse, error) {
	ctx.Logger().Info("new search term request", appcontext.Fields{"term": req.Term, "from": req.From, "to": req.To})

	ctx.Logger().Text("new domain model")
	domainTerm, err := domain.NewTerm(req.Term, req.From, req.To)
	if err != nil {
		ctx.Logger().Error("failed to create new domain term", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find in db first (caching)")
	term, err := h.termRepository.FindByTerm(ctx, req.Term, req.From)
	if err != nil {
		ctx.Logger().Error("failed to find in db first", err, appcontext.Fields{})
		return nil, err
	}
	if term != nil {
		ctx.Logger().Text("found in db, return")
		result := dto.SearchTermResponse{}.FromDomain(*term)
		return &result, nil
	}

	var (
		wg         sync.WaitGroup
		openaiData *domain.OpenAISearchResult
		scrapeData *domain.EnglishDictionaryScraperData
	)

	wg.Add(2)
	go func() {
		defer wg.Done()

		ctx.Logger().Text("search with Open AI")
		openaiData, err = h.openaiRepository.SearchTerm(ctx, req.Term, req.From, req.To)
		if err != nil {
			ctx.Logger().Error("failed to search with Open AI", err, appcontext.Fields{})
			return
		}
	}()

	go func() {
		defer wg.Done()

		ctx.Logger().Text("scrape English data")
		scrapeData, err = h.scraperRepository.GetEnglishDictionaryData(ctx, req.Term)
		if err != nil {
			ctx.Logger().Error("scrape English data", err, appcontext.Fields{})
			return
		}
	}()

	wg.Wait()

	if openaiData == nil {
		ctx.Logger().Error("term not found", nil, appcontext.Fields{})
		return nil, apperrors.Language.InvalidTerm
	}

	// assign term data
	_ = domainTerm.SetLanguage(openaiData.From.Language, openaiData.From.Definition, openaiData.From.Example)
	_ = domainTerm.SetLanguage(openaiData.To.Language, openaiData.To.Definition, openaiData.To.Example)
	domainTerm.SetSynonyms(openaiData.Synonyms)
	domainTerm.SetAntonyms(openaiData.Antonyms)
	if scrapeData != nil {
		domainTerm.SetLevel(scrapeData.Level)
		domainTerm.SetPartOfSpeech(scrapeData.PartOfSpeech)
		domainTerm.SetPhonetic(scrapeData.Phonetic)
		domainTerm.SetAudioURL(scrapeData.AudioURL)
	}

	ctx.Logger().Text("insert to database")
	if err = h.termRepository.CreateTerm(ctx, *domainTerm); err != nil {
		ctx.Logger().Error("failed to insert to database", err, appcontext.Fields{})
		return nil, err
	}

	// TODO: insert search history for performer

	ctx.Logger().Text("done search term request")
	result := dto.SearchTermResponse{}.FromDomain(*domainTerm)
	return &result, nil
}
