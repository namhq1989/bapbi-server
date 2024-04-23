package query

import (
	"sync"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"

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
	userHub                     domain.UserHub
}

func NewSearchTermHandler(
	termRepository domain.TermRepository,
	userSearchHistoryRepository domain.UserSearchHistoryRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
	userHub domain.UserHub,
) SearchTermHandler {
	return SearchTermHandler{
		termRepository:              termRepository,
		userSearchHistoryRepository: userSearchHistoryRepository,
		openaiRepository:            openaiRepository,
		scraperRepository:           scraperRepository,
		userHub:                     userHub,
	}
}

func (h SearchTermHandler) SearchTerm(ctx *appcontext.AppContext, performerID string, req dto.SearchTermRequest) (*dto.SearchTermResponse, error) {
	ctx.Logger().Info("new search term request", appcontext.Fields{"term": req.Term, "from": req.From, "to": req.To})

	ctx.Logger().Text("get user's subscription plan")
	plan, err := h.userHub.GetUserPlan(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to get user's subscription plan", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count total searched today")
	var (
		start = manipulation.StartOfToday()
		end   = time.Now()
	)
	totalSearched, err := h.userSearchHistoryRepository.CountTotalSearchedByTimeRange(ctx, performerID, start, end)
	if err != nil {
		ctx.Logger().Error("failed to count total searched today", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded := plan.IsExceededSearchLimitation(totalSearched); isExceeded {
		ctx.Logger().Error("exceeded search term limitation", nil, appcontext.Fields{"plan": plan.String(), "searched": totalSearched})
		return nil, apperrors.User.ExceededPlanLimitation
	}

	ctx.Logger().Info("still available to search term, create new domain model", appcontext.Fields{"searched": totalSearched})
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
		return &dto.SearchTermResponse{
			Term: dto.Term{}.FromDomain(*term, false),
		}, nil
	}

	ctx.Logger().Text("term not found in db, scrape from online dictionary")
	scrapeData, err := h.scraperRepository.GetEnglishDictionaryData(ctx, req.Term)
	if err != nil {
		ctx.Logger().Error("scrape Cambridge dictionary data failed", err, appcontext.Fields{})
	}
	if scrapeData == nil {
		ctx.Logger().Error("this term is an invalid vocabulary", nil, appcontext.Fields{})
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
		domainTerm.SetExamples(searchTermData.Examples)
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
	return &dto.SearchTermResponse{
		Term: dto.Term{}.FromDomain(*domainTerm, false),
	}, nil
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
