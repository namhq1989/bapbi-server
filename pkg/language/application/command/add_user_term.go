package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type AddUserTermHandler struct {
	termRepository     domain.TermRepository
	userTermRepository domain.UserTermRepository
	queueRepository    domain.QueueRepository
	userHub            domain.UserHub
	languageService    domain.LanguageService
}

func NewAddUserTermHandler(
	termRepository domain.TermRepository,
	userTermRepository domain.UserTermRepository,
	queueRepository domain.QueueRepository,
	userHub domain.UserHub,
	languageService domain.LanguageService,
) AddUserTermHandler {
	return AddUserTermHandler{
		termRepository:     termRepository,
		userTermRepository: userTermRepository,
		queueRepository:    queueRepository,
		userHub:            userHub,
		languageService:    languageService,
	}
}

func (h AddUserTermHandler) AddUserTerm(ctx *appcontext.AppContext, performerID, termID string, _ dto.AddUserTermRequest) (*dto.AddUserTermResponse, error) {
	ctx.Logger().Info("new add user term request", appcontext.Fields{"performer": performerID, "term": termID})

	ctx.Logger().Text("check today actions limitation")
	isExceeded, err := h.languageService.IsExceededAddTermLimitation(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to check today actions limitation", err, appcontext.Fields{})
		return nil, err
	}
	if isExceeded {
		ctx.Logger().Error("exceeded action limitation", nil, appcontext.Fields{})
		return nil, apperrors.User.ExceededPlanLimitation
	}

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
		return &dto.AddUserTermResponse{}, nil
	}

	ctx.Logger().Text("create new user term")
	userTerm, err := domain.NewUserTerm(performerID, term.ID, term.Term)
	if err != nil {
		ctx.Logger().Error("failed to create new user term", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("insert user term into database")
	if err = h.userTermRepository.AddUserTerm(ctx, *userTerm); err != nil {
		ctx.Logger().Error("failed to insert user term into database", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("add queue job")
	if err = h.queueRepository.EnqueueNewUserTerm(ctx, *userTerm); err != nil {
		ctx.Logger().Error("failed to add queue job", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done add user term request")
	return &dto.AddUserTermResponse{}, nil
}
