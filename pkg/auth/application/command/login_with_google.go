package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type LoginWithGoogleHandler struct {
	ssoRepository       domain.SSORepository
	authTokenRepository domain.AuthTokenRepository
	userHub             domain.UserHub
	jwtRepository       domain.JwtRepository
}

func NewLoginWithGoogleHandler(ssoRepository domain.SSORepository, authTokenRepository domain.AuthTokenRepository, userHub domain.UserHub, jwtRepository domain.JwtRepository) LoginWithGoogleHandler {
	return LoginWithGoogleHandler{
		ssoRepository:       ssoRepository,
		authTokenRepository: authTokenRepository,
		userHub:             userHub,
		jwtRepository:       jwtRepository,
	}
}

func (h LoginWithGoogleHandler) LoginWithGoogle(ctx *appcontext.AppContext, req dto.LoginWithGoogleRequest) (*dto.LoginWithGoogleResponse, error) {
	ctx.Logger().Info("new login with Google", appcontext.Fields{"token": req.Token})

	// get Google user data from token
	ctx.Logger().Info("get user data with Google token", appcontext.Fields{})
	googleUser, err := h.ssoRepository.GetUserDataWithGoogleToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	// find user by email
	ctx.Logger().Info("find user in database with email", appcontext.Fields{"email": googleUser.Email})
	user, err := h.userHub.GetOneByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, err
	}

	// generate tokens
	ctx.Logger().Info("user found, generate token", appcontext.Fields{})
	generatedTokens, err := h.jwtRepository.GenerateTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// persist refresh token
	ctx.Logger().Info("persist refresh token to database", appcontext.Fields{})
	err = h.authTokenRepository.CreateAuthToken(ctx, domain.RefreshToken{
		UserID: user.ID,
		Token:  generatedTokens.RefreshToken,
		Expiry: generatedTokens.RefreshTokenExpiry,
	})
	if err != nil {
		return nil, err
	}

	// response tokens
	ctx.Logger().Info("generate response's tokens data", appcontext.Fields{})
	tokens := &domain.Tokens{
		AccessToken:  generatedTokens.AccessToken,
		RefreshToken: generatedTokens.RefreshToken,
	}

	// return
	ctx.Logger().Info("done login with Google", appcontext.Fields{})
	return &dto.LoginWithGoogleResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
