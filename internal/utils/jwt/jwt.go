package appjwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	defaultAccessTokenTTL  = time.Minute * 15   // 15 minutes
	defaultRefreshTokenTTL = time.Hour * 24 * 7 // 1 week
)

type JWT struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type Result struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  time.Time
	RefreshTokenExpiry time.Time
}

func Init(accessTokenSecret, refreshTokenSecret string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) (*JWT, error) {
	// if ttl is zero, set default
	if accessTokenTTL.Seconds() == 0 {
		accessTokenTTL = defaultAccessTokenTTL
	}
	if refreshTokenTTL.Seconds() == 0 {
		refreshTokenTTL = defaultRefreshTokenTTL
	}

	return &JWT{
		accessTokenSecret:  []byte(accessTokenSecret),
		refreshTokenSecret: []byte(refreshTokenSecret),
		accessTokenTTL:     accessTokenTTL,
		refreshTokenTTL:    refreshTokenTTL,
	}, nil
}
