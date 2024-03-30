package query

import (
	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type GetAccessTokenByUserIdHandler struct {
	jwtRepository domain.JwtRepository
}

func NewGetAccessTokenByUserIdHandler(jwtRepository domain.JwtRepository) GetAccessTokenByUserIdHandler {
	return GetAccessTokenByUserIdHandler{
		jwtRepository: jwtRepository,
	}
}

func (h GetAccessTokenByUserIdHandler) GetAccessTokenByUserId(ctx *appcontext.AppContext, req dto.GetAccessTokenByUserIDRequest) (*dto.GetAccessTokenByUserIDResponse, error) {
	ctx.Logger().Info("get access token by user id", appcontext.Fields{"userID": req.UserID})

	if !database.IsValidObjectID(req.UserID) {
		ctx.Logger().Error("invalid user id", nil, appcontext.Fields{})
		return nil, apperrors.User.InvalidUserID
	}

	ctx.Logger().Info("generate new access token", appcontext.Fields{})
	token, err := h.jwtRepository.GenerateAccessToken(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Info("done get access token", appcontext.Fields{"token": token})
	return &dto.GetAccessTokenByUserIDResponse{AccessToken: token}, nil
}
