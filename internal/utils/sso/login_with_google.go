package sso

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

const googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"

type GoogleUserData struct {
	ID     string `json:"sub"`
	AUD    string `json:"aud"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"picture"`
}

func GetUserDataWithGoogleToken(ctx *appcontext.AppContext, clientID, token string) (*GoogleUserData, error) {
	// call api
	url := fmt.Sprintf(googleTokenInfoURL, token)
	resp, err := http.Get(fmt.Sprintf(googleTokenInfoURL, token))
	if err != nil {
		ctx.Logger().Error("fetch token info failed", err, appcontext.Fields{"url": url})
		return nil, err
	}

	// parse body
	defer func() { _ = resp.Body.Close() }()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.Logger().Error("parse response body failed", err, appcontext.Fields{})
		return nil, err
	}

	// map response to struct
	var result GoogleUserData
	if err = json.Unmarshal(contents, &result); err != nil {
		ctx.Logger().Error("map response to struct failed", err, appcontext.Fields{"contents": string(contents)})
		return nil, err
	}

	// verify that the token was issued to this application
	if result.AUD != clientID {
		err = errors.New("wrong audience")
		ctx.Logger().Error(err.Error(), err, appcontext.Fields{"response": result})
		return nil, err
	}

	// return
	return &result, nil
}
