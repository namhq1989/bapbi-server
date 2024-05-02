package service

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
)

func (s LanguageService) PersistUserActionHistory(ctx *appcontext.AppContext, userID, actionType string, data domain.UserActionHistoryData) error {
	ctx.Logger().Text("new user action history")

	action, err := domain.NewUserActionHistory(userID, actionType)
	if err != nil {
		ctx.Logger().Error("failed to create new user action history", err, appcontext.Fields{})
		return err
	}

	// set data
	action.SetData(data)

	ctx.Logger().Text("insert user action history to database")
	if err = s.userActionHistoryRepository.CreateUserActionHistory(ctx, *action); err != nil {
		ctx.Logger().Error("failed to insert user action history to database", err, appcontext.Fields{})
		return err
	}

	return nil
}
