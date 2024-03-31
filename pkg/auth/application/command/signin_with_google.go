package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type SignInWithGoogleHandler struct {
	ssoRepository       domain.SSORepository
	authTokenRepository domain.AuthTokenRepository
	userHub             domain.UserHub
	jwtRepository       domain.JwtRepository
}

func NewSignInWithGoogleHandler(ssoRepository domain.SSORepository, authTokenRepository domain.AuthTokenRepository, userHub domain.UserHub, jwtRepository domain.JwtRepository) SignInWithGoogleHandler {
	return SignInWithGoogleHandler{
		ssoRepository:       ssoRepository,
		authTokenRepository: authTokenRepository,
		userHub:             userHub,
		jwtRepository:       jwtRepository,
	}
}

func (h SignInWithGoogleHandler) SignInWithGoogle(ctx *appcontext.AppContext, req dto.SignInWithGoogleRequest) (*dto.SignInWithGoogleResponse, error) {
	ctx.Logger().Info("new sign in with Google", appcontext.Fields{"token": req.Token})

	// get Google user data from token
	ctx.Logger().Info("get user data with Google token", appcontext.Fields{})
	googleUser, err := h.ssoRepository.GetUserDataWithGoogleToken(ctx, req.Token)
	if err != nil {
		ctx.Logger().Error("failed to get user data with Google token", err, appcontext.Fields{"token": req.Token})
		return nil, err
	}

	// find user by email
	ctx.Logger().Info("find user in database with email", appcontext.Fields{"email": googleUser.Email})
	user, err := h.userHub.GetOneByEmail(ctx, googleUser.Email)
	if err != nil {
		ctx.Logger().Error("failed to get user by email", err, appcontext.Fields{"email": googleUser.Email})
		return nil, err
	}
	if user == nil {
		ctx.Logger().Info("user not found", appcontext.Fields{"email": googleUser.Email})
		return nil, apperrors.User.UserNotFound
	}

	// generate tokens
	ctx.Logger().Info("user found, generate token", appcontext.Fields{})
	generatedTokens, err := h.jwtRepository.GenerateTokens(ctx, user.ID)
	if err != nil {
		ctx.Logger().Error("failed to generate tokens", err, appcontext.Fields{"userId": user.ID})
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
		ctx.Logger().Error("failed to persist refresh token", err, appcontext.Fields{"userId": user.ID})
		return nil, err
	}

	// response tokens
	ctx.Logger().Info("generate response's tokens data", appcontext.Fields{})
	tokens := &domain.Tokens{
		AccessToken:  generatedTokens.AccessToken,
		RefreshToken: generatedTokens.RefreshToken,
	}

	// return
	ctx.Logger().Info("done sign in with Google", appcontext.Fields{})
	return &dto.SignInWithGoogleResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
