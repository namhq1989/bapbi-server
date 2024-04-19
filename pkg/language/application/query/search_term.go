package query

import (
	"sync"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SearchTermHandler struct {
	termRepository              domain.TermRepository
	userSearchHistoryRepository domain.UserSearchHistoryRepository
	openaiRepository            domain.OpenAIRepository
	scraperRepository           domain.ScraperRepository
}

func NewSearchTermHandler(
	termRepository domain.TermRepository,
	userSearchHistoryRepository domain.UserSearchHistoryRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
) SearchTermHandler {
	return SearchTermHandler{
		termRepository:              termRepository,
		userSearchHistoryRepository: userSearchHistoryRepository,
		openaiRepository:            openaiRepository,
		scraperRepository:           scraperRepository,
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
		ctx.Logger().Text("found in db")

		if err = h.insertUserSearchHistory(ctx, performerID, req.Term, true); err != nil {
			return nil, err
		}

		ctx.Logger().Text("respond")
		result := dto.SearchTermResponse{}.FromDomain(*term)
		return &result, nil
	}

	ctx.Logger().Text("term not found in db, scrape from online dictionary")
	scrapeData, err := h.scraperRepository.GetEnglishDictionaryData(ctx, req.Term)
	if err != nil {
		ctx.Logger().Error("scrape Cambridge dictionary data failed", err, appcontext.Fields{})
	}
	if scrapeData == nil {
		ctx.Logger().Text("this term is an invalid vocabulary")
		if err = h.insertUserSearchHistory(ctx, performerID, req.Term, false); err != nil {
			return nil, err
		}
		return nil, apperrors.Language.InvalidTerm
	}

	var (
		wg                    sync.WaitGroup
		searchTermData        *domain.OpenAISearchTermResult
		semanticRelationsData *domain.OpenAISearchSemanticRelationsResult
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		ctx.Logger().Text("search semantic relations with Open AI")
		semanticRelationsData, err = h.openaiRepository.SearchSemanticRelations(ctx, req.Term, req.From)
		if err != nil {
			ctx.Logger().Error("failed to search semantic relations with Open AI", err, appcontext.Fields{})
		}
	}()

	go func() {
		defer wg.Done()

		// determine the term is valid or not
		ctx.Logger().Text("search term with Open AI")
		searchTermData, err = h.openaiRepository.SearchTerm(ctx, req.Term, req.From, req.To)
		if err != nil {
			ctx.Logger().Error("failed to search term with Open AI", err, appcontext.Fields{})
		}
	}()

	wg.Wait()

	// assign term data
	if searchTermData != nil {
		_ = domainTerm.SetLanguage(searchTermData.From.Language, searchTermData.From.Definition, searchTermData.From.Example)
		_ = domainTerm.SetLanguage(searchTermData.To.Language, searchTermData.To.Definition, searchTermData.To.Example)
	}

	if semanticRelationsData != nil {
		domainTerm.SetSynonyms(semanticRelationsData.Synonyms)
		domainTerm.SetAntonyms(semanticRelationsData.Antonyms)
	}

	domainTerm.SetLevel(scrapeData.Level)
	domainTerm.SetPartOfSpeech(scrapeData.PartOfSpeech)
	domainTerm.SetPhonetic(scrapeData.Phonetic)
	domainTerm.SetAudioURL(scrapeData.AudioURL)
	domainTerm.SetReferenceURL(scrapeData.ReferenceURL)

	ctx.Logger().Text("insert to database")
	if err = h.termRepository.CreateTerm(ctx, *domainTerm); err != nil {
		ctx.Logger().Error("failed to insert to database", err, appcontext.Fields{})
		return nil, err
	}

	if err = h.insertUserSearchHistory(ctx, performerID, req.Term, true); err != nil {
		return nil, err
	}

	ctx.Logger().Text("done search term request")
	result := dto.SearchTermResponse{}.FromDomain(*domainTerm)
	return &result, nil
}

func (h SearchTermHandler) insertUserSearchHistory(ctx *appcontext.AppContext, performerID, term string, isValid bool) error {
	ctx.Logger().Text("new user search history")

	domainHistory, err := domain.NewUserSearchHistory(performerID, term, isValid)
	if err != nil {
		ctx.Logger().Error("failed to create new user search history", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("insert user search history to database")
	if err = h.userSearchHistoryRepository.CreateUserSearchHistory(ctx, *domainHistory); err != nil {
		ctx.Logger().Error("failed to insert user search history to database", err, appcontext.Fields{})
		return err
	}

	return nil
}
