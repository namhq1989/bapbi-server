package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

func (w Workers) FeaturedWord(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx = appcontext.New(bgCtx)

		fromLanguage = domain.LanguageEnglish.String()
		toLanguage   = domain.LanguageVietnamese.String()
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	// call OpenAI to get featured word
	ctx.Logger().Text("call OpenAI to get featured word")
	featureWordData, err := w.openaiRepository.FeaturedWord(ctx, fromLanguage)
	if err != nil {
		ctx.Logger().Error("failed to call OpenAI to get featured word", err, appcontext.Fields{})
		return err
	}

	// check word already existed or not
	ctx.Logger().Info("check word already existed or not", appcontext.Fields{"word": featureWordData.Word})
	term, err := w.termRepository.FindByTerm(ctx, featureWordData.Word, fromLanguage)
	if err != nil {
		ctx.Logger().Error("failed to find word in db", nil, appcontext.Fields{})
		return err
	}
	if term != nil {
		ctx.Logger().Text("word already existed")
		return errors.New("word already existed")
	}

	ctx.Logger().Text("scrape from online dictionary")
	scrapeData, err := w.scraperRepository.GetEnglishDictionaryData(ctx, featureWordData.Word)
	if err != nil {
		ctx.Logger().Error("scrape Cambridge dictionary data failed", err, appcontext.Fields{})
	}
	if scrapeData == nil {
		ctx.Logger().Error("this term is an invalid vocabulary", nil, appcontext.Fields{})
		return apperrors.Language.InvalidTerm
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
		semanticRelationsData, err = w.openaiRepository.SearchSemanticRelations(ctx, featureWordData.Word, fromLanguage)
		if err != nil {
			ctx.Logger().Error("failed to search semantic relations with Open AI", err, appcontext.Fields{})
		}
	}()

	go func() {
		defer wg.Done()

		ctx.Logger().Text("search term with Open AI")
		searchTermData, err = w.openaiRepository.SearchTerm(ctx, featureWordData.Word, fromLanguage, toLanguage)
		if err != nil {
			ctx.Logger().Error("failed to search term with Open AI", err, appcontext.Fields{})
		}
	}()

	wg.Wait()

	// assign data
	domainTerm, _ := domain.NewTerm(featureWordData.Word, fromLanguage, toLanguage)

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
	domainTerm.SetIsFeatured(true)

	ctx.Logger().Text("insert to database")
	if err = w.termRepository.CreateTerm(ctx, *domainTerm); err != nil {
		ctx.Logger().Error("failed to insert to database", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}
