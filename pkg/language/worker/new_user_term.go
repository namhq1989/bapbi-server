package worker

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

func (w Workers) NewUserTerm(bgCtx context.Context, t *asynq.Task) error {
	var (
		ctx      = appcontext.New(bgCtx)
		userTerm domain.UserTerm
	)

	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Info("unmarshal task payload", appcontext.Fields{})
	if err := json.Unmarshal(t.Payload(), &userTerm); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	//
	// TASK
	//

	if err := w.addUserVocabularyExercise(ctx, userTerm); err != nil {
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}

func (w Workers) addUserVocabularyExercise(ctx *appcontext.AppContext, userTerm domain.UserTerm) error {
	ctx.Logger().Text("add user vocabulary exercise")

	ctx.Logger().Text("check user vocabulary exercise existence")
	isExisted, err := w.userVocabularyExerciseRepository.IsExerciseCreated(ctx, userTerm.UserID, userTerm.TermID)
	if err != nil {
		ctx.Logger().Error("failed to check user vocabulary exercise", err, appcontext.Fields{})
		return err
	}
	if isExisted {
		ctx.Logger().Text("user vocabulary exercise already created, stop the flow")
		return nil
	}

	ctx.Logger().Text("random English grammar tense")
	tenseCode := domain.RandomGrammarTenseCode()

	ctx.Logger().Text("create user vocabulary exercise")
	exercise, err := domain.NewUserVocabularyExercise(userTerm.UserID, userTerm.TermID, userTerm.Term, domain.LanguageEnglish.String(), tenseCode.String())
	if err != nil {
		ctx.Logger().Error("failed to create user vocabulary exercise", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("persist user vocabulary exercise to database")
	if err = w.userVocabularyExerciseRepository.CreateUserVocabularyExercise(ctx, *exercise); err != nil {
		ctx.Logger().Error("failed to persist user vocabulary exercise", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("done add user vocabulary exercise")
	return nil
}
