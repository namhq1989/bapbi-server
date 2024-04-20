package command

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type AddTermHandler struct {
	termRepository     domain.TermRepository
	userTermRepository domain.UserTermRepository
	userHub            domain.UserHub
}

func NewAddTermHandler(termRepository domain.TermRepository, userTermRepository domain.UserTermRepository, userHub domain.UserHub) AddTermHandler {
	return AddTermHandler{
		termRepository:     termRepository,
		userTermRepository: userTermRepository,
		userHub:            userHub,
	}
}

func (h AddTermHandler) AddTerm(ctx *appcontext.AppContext, performerID, termID string, _ dto.AddTermRequest) (*dto.AddTermResponse, error) {
	ctx.Logger().Info("new add term request", appcontext.Fields{"performer": performerID, "term": termID})

	ctx.Logger().Text("find term in db")
	term, err := h.termRepository.FindByID(ctx, termID)
	if err != nil {
		ctx.Logger().Error("failed to find term in db", err, appcontext.Fields{})
		return nil, err
	}
	if term == nil {
		ctx.Logger().Error("term not found", nil, appcontext.Fields{})
		return nil, apperrors.Language.TermNotFound
	}

	ctx.Logger().Text("check if user already added this term or not")
	isAdded, err := h.userTermRepository.IsUserTermAdded(ctx, performerID, term.Term)
	if err != nil {
		ctx.Logger().Error("failed to check if user already added this term or not", err, appcontext.Fields{})
		return nil, err
	}
	if isAdded {
		ctx.Logger().Text("user already added this term, respond")
		return &dto.AddTermResponse{}, nil
	}

	ctx.Logger().Text("get user's subscription plan")
	plan, err := h.userHub.GetUserPlan(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to get user's subscription plan", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count total terms added today")
	var (
		start = manipulation.StartOfToday()
		end   = time.Now()
	)
	totalAdded, err := h.userTermRepository.CountTotalTermAddedByTimeRange(ctx, performerID, start, end)
	if err != nil {
		ctx.Logger().Error("failed to count total terms added today", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded := plan.IsExceededAddTermLimitation(totalAdded); isExceeded {
		ctx.Logger().Error("exceeded add term limitation", nil, appcontext.Fields{})
		return nil, apperrors.User.ExceededPlanLimitation
	}

	ctx.Logger().Info("still available to add term, create new user term", appcontext.Fields{"added": totalAdded})
	userTerm, err := domain.NewUserTerm(performerID, term.Term)
	if err != nil {
		ctx.Logger().Error("failed to create new user term", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("insert user term into database")
	if err = h.userTermRepository.AddUserTerm(ctx, *userTerm); err != nil {
		ctx.Logger().Error("failed to insert user term into database", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("add user's term successfully")
	return &dto.AddTermResponse{}, nil
}
