package domain

import "github.com/namhq1989/bapbi-server/internal/utils/appcontext"

type LanguageService interface {
	IsExceededAddTermLimitation(ctx *appcontext.AppContext, userID string) (bool, error)
	IsExceededActionLimitation(ctx *appcontext.AppContext, userID string) (bool, error)
	PersistUserActionHistory(ctx *appcontext.AppContext, userID, actionType string, data UserActionHistoryData) error
}
