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
	userHub           domain.UserHub
	languageService   domain.LanguageService
}

func NewSearchTermHandler(
	termRepository domain.TermRepository,
	openaiRepository domain.OpenAIRepository,
	scraperRepository domain.ScraperRepository,
	userHub domain.UserHub,
	languageService domain.LanguageService,
) SearchTermHandler {
	return SearchTermHandler{
		termRepository:    termRepository,
		openaiRepository:  openaiRepository,
		scraperRepository: scraperRepository,
		userHub:           userHub,
		languageService:   languageService,
	}
}

func (h SearchTermHandler) SearchTerm(ctx *appcontext.AppContext, performerID string, req dto.SearchTermRequest) (*dto.SearchTermResponse, error) {
	ctx.Logger().Info("new search term request", appcontext.Fields{"term": req.Term, "from": req.From, "to": req.To})

	ctx.Logger().Text("check today actions limitation")
	isExceeded, err := h.languageService.IsExceededActionLimitation(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to check today actions limitation", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded {
		ctx.Logger().Error("exceeded action limitation", nil, appcontext.Fields{})
		return nil, apperrors.User.ExceededPlanLimitation
	}

	ctx.Logger().Text("create new domain model")
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

		if err = h.languageService.PersistUserActionHistory(ctx, performerID, domain.UserActionTypeSearchTerm.String(), domain.UserActionHistoryData{
			Term:    req.Term,
			IsValid: true,
		}); err != nil {
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
		if err = h.languageService.PersistUserActionHistory(ctx, performerID, domain.UserActionTypeSearchTerm.String(), domain.UserActionHistoryData{
			Term:    req.Term,
			IsValid: false,
		}); err != nil {
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

	if err = h.languageService.PersistUserActionHistory(ctx, performerID, domain.UserActionTypeSearchTerm.String(), domain.UserActionHistoryData{
		Term:    req.Term,
		IsValid: true,
	}); err != nil {
		return nil, err
	}

	ctx.Logger().Text("done search term request")
	return &dto.SearchTermResponse{
		Term: dto.Term{}.FromDomain(*domainTerm, false),
	}, nil
}
