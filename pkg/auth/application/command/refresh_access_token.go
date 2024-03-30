package command

import (
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type RefreshAccessTokenHandler struct {
	authTokenRepository domain.AuthTokenRepository
	userHub             domain.UserHub
	jwtRepository       domain.JwtRepository
}

func NewRefreshAccessTokenHandler(authTokenRepository domain.AuthTokenRepository, userHub domain.UserHub, jwtRepository domain.JwtRepository) RefreshAccessTokenHandler {
	return RefreshAccessTokenHandler{
		authTokenRepository: authTokenRepository,
		userHub:             userHub,
		jwtRepository:       jwtRepository,
	}
}

func (h RefreshAccessTokenHandler) RefreshAccessToken(ctx *appcontext.AppContext, req dto.RefreshAccessTokenRequest) (*dto.RefreshAccessTokenResponse, error) {
	ctx.Logger().Info("refresh access token", appcontext.Fields{"userID": req.UserID, "refreshToken": req.RefreshToken})

	// find user first
	ctx.Logger().Info("find user in database", appcontext.Fields{})
	user, err := h.userHub.GetOneByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// fetch auth token document
	ctx.Logger().Info("find auth token in database", appcontext.Fields{})
	authToken, err := h.authTokenRepository.FindAuthToken(ctx, req.UserID, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check expiration time
	ctx.Logger().Info("check expiration time", appcontext.Fields{"expiry": authToken.Expiry.String()})
	if authToken.Expiry.Before(time.Now()) {
		// delete token
		defer func() { _ = h.authTokenRepository.DeleteAuthToken(ctx, *authToken) }()
		return nil, apperrors.Auth.InvalidAuthToken
	}

	// generate new access token
	ctx.Logger().Info("generate new access token", appcontext.Fields{})
	if accessToken, err := h.jwtRepository.GenerateAccessToken(ctx, user.ID); err != nil {
		return nil, err
	} else {
		ctx.Logger().Info("done refresh access token", appcontext.Fields{})
		return &dto.RefreshAccessTokenResponse{AccessToken: accessToken}, nil
	}
}
