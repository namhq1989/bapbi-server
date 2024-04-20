package pagetoken

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// PageToken ...
type PageToken struct {
	Page      int64
	Timestamp time.Time
}

func getDefaultPageToken() PageToken {
	return PageToken{
		Page:      0,
		Timestamp: time.Now(),
	}
}

func encode(page int64, timestamp time.Time) string {
	tokenData := PageToken{
		Page:      page,
		Timestamp: timestamp,
	}
	tokenString, _ := json.Marshal(tokenData)
	encodedString := base64.StdEncoding.EncodeToString(tokenString)
	return encodedString
}

func Decode(token string) PageToken {
	if token == "" {
		return getDefaultPageToken()
	}

	// Decode string
	decoded, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		return getDefaultPageToken()
	}

	// Parse token
	var pageToken PageToken
	err = json.Unmarshal(decoded, &pageToken)
	if err != nil {
		return getDefaultPageToken()
	}

	return pageToken
}

func NewWithPage(page int64) string {
	return encode(page, time.Now())
}

func NewWithTimestamp(timestamp time.Time) string {
	return encode(0, timestamp)
}
