// Package api ...
package api

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

type Token struct {
	AccessToken string
	ExpirationTime time.Time
	RefreshToken string
}

// <---------------------------------------------------------------------------------------------------->



func (token *Token) GetAccessToken() string {
	currentTime := time.Now()

	if (token.ExpirationTime.Before(currentTime)) {
		token.refreshAccessToken(currentTime)
	}

	return token.AccessToken
}


func (token *Token) refreshAccessToken(currentTime time.Time) error {
    parameters := map[string]string{
        "grant_type": "refresh_token",
        "refresh_token" : token.RefreshToken,
    }
    headers := map[string]string{
        "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)),
        "Content-Type" : "application/x-www-form-urlencoded",
    }

	responseMap, err := util.MakePOSTRequest("https://accounts.spotify.com/api/token", parameters, headers)
	if err != nil {
        return err
    }

	token.AccessToken = responseMap["access_token"].(string)
	token.ExpirationTime = currentTime.Add(time.Hour)
	token.RefreshToken = responseMap["refresh_token"].(string)

	return nil
}