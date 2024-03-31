package command

import (
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/dto"
)

type SignUpWithGoogleHandler struct {
	ssoRepository       domain.SSORepository
	authTokenRepository domain.AuthTokenRepository
	userHub             domain.UserHub
	jwtRepository       domain.JwtRepository
}

func NewSignUpWithGoogleHandler(ssoRepository domain.SSORepository, authTokenRepository domain.AuthTokenRepository, userHub domain.UserHub, jwtRepository domain.JwtRepository) SignUpWithGoogleHandler {
	return SignUpWithGoogleHandler{
		ssoRepository:       ssoRepository,
		authTokenRepository: authTokenRepository,
		userHub:             userHub,
		jwtRepository:       jwtRepository,
	}
}

func (h SignUpWithGoogleHandler) SignUpWithGoogle(ctx *appcontext.AppContext, req dto.SignUpWithGoogleRequest) (*dto.SignUpWithGoogleResponse, error) {
	ctx.Logger().Info("new sign up with Google", appcontext.Fields{"token": req.Token})

	// get Google user data from token
	ctx.Logger().Info("get user data with Google token", appcontext.Fields{})
	googleUser, err := h.ssoRepository.GetUserDataWithGoogleToken(ctx, req.Token)
	if err != nil {
		ctx.Logger().Error("failed to get user data with Google token", err, appcontext.Fields{"token": req.Token})
		return nil, err
	}

	// find user by email
	ctx.Logger().Info("find user in database with email", appcontext.Fields{"email": googleUser.Email})
	dbUser, err := h.userHub.GetOneByEmail(ctx, googleUser.Email)
	if err != nil {
		ctx.Logger().Error("failed to get user by email", err, appcontext.Fields{"email": googleUser.Email})
		return nil, err
	}
	if dbUser != nil {
		ctx.Logger().Info("user already existed with this email", appcontext.Fields{"email": googleUser.Email})
		return nil, apperrors.User.GoogleEmailAlreadyRegistered
	}

	// create user
	ctx.Logger().Info("create user in database", appcontext.Fields{})
	userID, err := h.userHub.CreateUser(ctx, domain.User{
		Name:  googleUser.Name,
		Email: googleUser.Email,
	})
	if err != nil {
		ctx.Logger().Error("failed to create user", err, appcontext.Fields{"email": googleUser.Email})
		return nil, err
	}

	// generate tokens
	ctx.Logger().Info("generate token", appcontext.Fields{})
	generatedTokens, err := h.jwtRepository.GenerateTokens(ctx, userID)
	if err != nil {
		ctx.Logger().Error("failed to generate tokens", err, appcontext.Fields{"userId": userID})
		return nil, err
	}

	// persist refresh token
	ctx.Logger().Info("persist refresh token to database", appcontext.Fields{})
	err = h.authTokenRepository.CreateAuthToken(ctx, domain.RefreshToken{
		UserID: userID,
		Token:  generatedTokens.RefreshToken,
		Expiry: generatedTokens.RefreshTokenExpiry,
	})
	if err != nil {
		ctx.Logger().Error("failed to persist refresh token", err, appcontext.Fields{"userId": userID})
		return nil, err
	}

	// response tokens
	ctx.Logger().Info("generate response's tokens data", appcontext.Fields{})
	tokens := &domain.Tokens{
		AccessToken:  generatedTokens.AccessToken,
		RefreshToken: generatedTokens.RefreshToken,
	}

	// return
	ctx.Logger().Info("done sign up with Google", appcontext.Fields{})
	return &dto.SignUpWithGoogleResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
