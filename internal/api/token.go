// Package api ...
package api

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// Token is a type to hold our token data
type Token struct {
	accessToken string
	expirationTime time.Time
	refreshToken string
}

// <---------------------------------------------------------------------------------------------------->



// GetAccessToken is a getter for the access token that always ensures it's up to date
func (token *Token) getAccessToken() string {
	currentTime := time.Now()

	// chech if the access token is still usable
	if (token.expirationTime.Before(currentTime)) {
		err := token.refreshAccessToken(currentTime)
		if err != nil {
			util.LogError(err)
		}
	}

	return token.accessToken
}



func (token *Token) GetAccessTokenHeader() map[string]string {
	return map[string]string{"Authorization": "Bearer " + token.getAccessToken(),}
}



// refreshAccessToken uses the refreshToken to exchange for a new accessToken
func (token *Token) refreshAccessToken(currentTime time.Time) error {
	parameters := map[string]string{
		"grant_type": "refresh_token",
		"refresh_token" : token.refreshToken,
	}
	headers := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
		"Content-Type" : "application/x-www-form-urlencoded",
	}

	responseMap, err := util.MakeHTTPRequest("POST", tokenURL, headers, parameters, nil)
	if err != nil {
		return errors.New("couldn't POST request refreshed token: " + err.Error())
	}

	token.accessToken = responseMap["access_token"].(string)
	token.expirationTime = time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))) * time.Second)

	return nil
}