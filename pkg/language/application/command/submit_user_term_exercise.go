package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/httprespond"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type SubmitUserTermExerciseHandler struct {
	userTermRepository         domain.UserTermRepository
	UserTermExerciseRepository domain.UserTermExerciseRepository
	openaiRepository           domain.OpenAIRepository
	languageService            domain.LanguageService
}

func NewSubmitUserTermExerciseHandler(
	userTermRepository domain.UserTermRepository,
	UserTermExerciseRepository domain.UserTermExerciseRepository,
	openaiRepository domain.OpenAIRepository,
	languageService domain.LanguageService,
) SubmitUserTermExerciseHandler {
	return SubmitUserTermExerciseHandler{
		userTermRepository:         userTermRepository,
		UserTermExerciseRepository: UserTermExerciseRepository,
		openaiRepository:           openaiRepository,
		languageService:            languageService,
	}
}

func (h SubmitUserTermExerciseHandler) SubmitUserTermExercise(ctx *appcontext.AppContext, performerID string, req dto.SubmitUserTermExerciseRequest) (*dto.SubmitUserTermExerciseResponse, error) {
	ctx.Logger().Info("new submit user term exercise request", appcontext.Fields{"performer": performerID, "exerciseId": req.ExerciseID, "content": req.Content})

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
	exercise, err := h.UserTermExerciseRepository.FindByExerciseID(ctx, req.ExerciseID)
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
	assessment, err := h.openaiRepository.AssessTermExercise(ctx, exercise.Language.String(), exercise.Term, exercise.Tense.String(), req.Content)
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

	ctx.Logger().Text("get and set status based on assessment")
	status := exercise.GetStatusBasedOnAssessment()
	exercise.SetStatus(status)

	ctx.Logger().Text("update user term exercise in db")
	err = h.UserTermExerciseRepository.UpdateUserTermExercise(ctx, *exercise)
	if err != nil {
		ctx.Logger().Error("failed to update user term exercise in db", err, appcontext.Fields{})
		return nil, err
	}

	dtoAssessment := dto.UserTermExerciseAssessment{}.FromDomain(exercise.Assessment)
	return &dto.SubmitUserTermExerciseResponse{
		Status:      exercise.Status.String(),
		CompletedAt: httprespond.NewTimeResponse(exercise.CompletedAt),
		Assessment:  *dtoAssessment,
	}, nil
}
