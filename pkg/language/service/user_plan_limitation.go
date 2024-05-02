package service

import (
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
)

func (s LanguageService) IsExceededActionLimitation(ctx *appcontext.AppContext, userID string) (bool, error) {
	ctx.Logger().Text("get user's subscription plan")
	plan, err := s.userHub.GetUserPlan(ctx, userID)
	if err != nil {
		ctx.Logger().Error("failed to get user's subscription plan", err, appcontext.Fields{})
		return false, err
	}

	ctx.Logger().Text("count total actions today")
	var (
		start = manipulation.StartOfToday()
		end   = time.Now()
	)
	totalActions, err := s.userActionHistoryRepository.CountTotalActionsByTimeRange(ctx, userID, start, end)
	if err != nil {
		ctx.Logger().Error("failed to count total actions today", err, appcontext.Fields{})
		return false, err
	}
	if isExceeded := plan.IsExceededActionLimitation(totalActions); isExceeded {
		ctx.Logger().Error("exceeded action limitation", nil, appcontext.Fields{"plan": plan.String(), "actions": totalActions})
		return true, nil
	}

	ctx.Logger().Info("still available to perform an action", appcontext.Fields{"actions": totalActions})
	return false, nil
}

func (s LanguageService) IsExceededAddTermLimitation(ctx *appcontext.AppContext, userID string) (bool, error) {
	ctx.Logger().Text("get user's subscription plan")
	plan, err := s.userHub.GetUserPlan(ctx, userID)
	if err != nil {
		ctx.Logger().Error("failed to get user's subscription plan", err, appcontext.Fields{})
		return false, err
	}

	ctx.Logger().Text("count total terms added today")
	var (
		start = manipulation.StartOfToday()
		end   = time.Now()
	)
	totalAdded, err := s.userTermRepository.CountTotalTermAddedByTimeRange(ctx, userID, start, end)
	if err != nil {
		ctx.Logger().Error("failed to count total terms added today", err, appcontext.Fields{})
		return false, err
	}

	if isExceeded := plan.IsExceededAddTermLimitation(totalAdded); isExceeded {
		ctx.Logger().Error("exceeded add term limitation", nil, appcontext.Fields{"plan": plan.String(), "terms added": totalAdded})
		return false, apperrors.User.ExceededPlanLimitation
	}

	ctx.Logger().Info("still available to perform an action", appcontext.Fields{"terms added": totalAdded})
	return false, nil
}
