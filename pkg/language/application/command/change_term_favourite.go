package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/dto"
)

type ChangeTermFavouriteHandler struct {
	userTermRepository domain.UserTermRepository
}

func NewChangeTermFavouriteHandler(
	userTermRepository domain.UserTermRepository,
) ChangeTermFavouriteHandler {
	return ChangeTermFavouriteHandler{
		userTermRepository: userTermRepository,
	}
}

func (h ChangeTermFavouriteHandler) ChangeTermFavourite(ctx *appcontext.AppContext, performerID, userTermID string, _ dto.ChangeTermFavouriteRequest) (*dto.ChangeTermFavouriteResponse, error) {
	ctx.Logger().Info("new change user term favourite request", appcontext.Fields{"performer": performerID, "userTerm": userTermID})

	ctx.Logger().Text("find user term in db")
	term, err := h.userTermRepository.FindUserTermByID(ctx, userTermID)
	if err != nil {
		ctx.Logger().Error("failed to find user term in db", err, appcontext.Fields{})
		return nil, err
	}
	if term == nil || term.UserID != performerID {
		ctx.Logger().Error("user term not found", nil, appcontext.Fields{})
		return nil, apperrors.Language.TermNotFound
	}

	newValue := !term.IsFavourite
	ctx.Logger().Info("set new favourite value", appcontext.Fields{"isFavourite": newValue})
	term.SetIsFavourite(newValue)

	ctx.Logger().Text("update user term in db")
	if err = h.userTermRepository.UpdateUserTerm(ctx, *term); err != nil {
		ctx.Logger().Error("failed to update user term in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done change term favourite")
	return &dto.ChangeTermFavouriteResponse{
		IsFavourite: newValue,
	}, nil
}
