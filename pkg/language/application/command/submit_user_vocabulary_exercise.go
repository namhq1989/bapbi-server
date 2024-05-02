package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SubmitUserVocabularyExerciseHandler struct {
	userTermRepository               domain.UserTermRepository
	userVocabularyExerciseRepository domain.UserVocabularyExerciseRepository
	openaiRepository                 domain.OpenAIRepository
	languageService                  domain.LanguageService
}

func NewSubmitUserVocabularyExerciseHandler(
	userTermRepository domain.UserTermRepository,
	userVocabularyExerciseRepository domain.UserVocabularyExerciseRepository,
	openaiRepository domain.OpenAIRepository,
	languageService domain.LanguageService,
) SubmitUserVocabularyExerciseHandler {
	return SubmitUserVocabularyExerciseHandler{
		userTermRepository:               userTermRepository,
		userVocabularyExerciseRepository: userVocabularyExerciseRepository,
		openaiRepository:                 openaiRepository,
		languageService:                  languageService,
	}
}

func (h SubmitUserVocabularyExerciseHandler) SubmitUserVocabularyExercise(ctx *appcontext.AppContext, performerID string, req dto.SubmitUserVocabularyExerciseRequest) (*dto.SubmitUserVocabularyExerciseResponse, error) {
	ctx.Logger().Info("new submit user vocabulary exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseID, "content": req.Content})

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

	ctx.Logger().Text("find exercise in db")
	exercise, err := h.userVocabularyExerciseRepository.FindByExerciseID(ctx, req.ExerciseID)
	if err != nil {
		ctx.Logger().Error("failed to find exercise in db", err, appcontext.Fields{})
		return nil, err
	}
	if exercise == nil || !exercise.IsOwner(performerID) {
		ctx.Logger().Error("exercise not found, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseNotFound
	}
	if exercise.IsCompleted() {
		ctx.Logger().Error("exercise already completed, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Language.ExerciseAlreadyCompleted
	}

	ctx.Logger().Text("call OpenAI's api to assess the exercise's content")
	assessment, err := h.openaiRepository.AssessVocabularyExercise(ctx, exercise.Language.String(), exercise.Term, exercise.Tense.String(), req.Content)
	if err != nil {
		ctx.Logger().Error("failed to call OpenAI's api to assess the exercise's content", err, appcontext.Fields{})
		return nil, err
	}
	_ = h.languageService.PersistUserActionHistory(ctx, performerID, domain.UserActionTypeSubmitVocabularyExercise.String(), domain.UserActionHistoryData{ExerciseID: exercise.ID})
	if assessment == nil {
		ctx.Logger().Error("failed to assess the exercise's content, respond error", nil, appcontext.Fields{})
		return nil, apperrors.Common.BadRequest
	}

	ctx.Logger().Text("set content")
	if err = exercise.SetContent(req.Content); err != nil {
		ctx.Logger().Error("failed to set content", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("set assessment")
	exercise.SetAssessment(assessment.IsVocabularyCorrect, assessment.VocabularyIssue, assessment.IsTenseCorrect, assessment.TenseIssue, assessment.GrammarIssues, assessment.ImprovementSuggestions)

	ctx.Logger().Text("set user vocabulary exercise status to completed")
	exercise.SetCompleted()

	ctx.Logger().Text("update user vocabulary exercise in db")
	err = h.userVocabularyExerciseRepository.UpdateUserVocabularyExercise(ctx, *exercise)
	if err != nil {
		ctx.Logger().Error("failed to update user vocabulary exercise in db", err, appcontext.Fields{})
		return nil, err
	}

	dtoAssessment := dto.UserVocabularyExerciseAssessment{}.FromDomain(*exercise.Assessment)
	return &dto.SubmitUserVocabularyExerciseResponse{
		Status:      exercise.Status.String(),
		CompletedAt: httprespond.NewTimeResponse(exercise.CompletedAt),
		Assessment:  dtoAssessment,
	}, nil
}
