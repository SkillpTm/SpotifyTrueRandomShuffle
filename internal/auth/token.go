// Package auth is responsible for handling everything to do with the authorization and token exchanges with Spotify
package auth

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

// Token is a type to hold our token data
type Token struct {
	accessToken    string
	expirationTime time.Time
	refreshToken   string
}

// <---------------------------------------------------------------------------------------------------->

// GetAccessTokenHeader generates the auth header needed fro essentially all API calls with a token
func (token *Token) GetAccessTokenHeader() map[string]string {
	return map[string]string{"Authorization": "Bearer " + token.getAccessToken()}
}

// GetAccessToken is a getter for the access token that always ensures it's up to date
func (token *Token) getAccessToken() string {
	currentTime := time.Now()

	// chech if the access token is still usable
	if token.expirationTime.Before(currentTime) {
		err := token.refreshAccessToken()
		if err != nil {
			util.LogError(fmt.Errorf("couldn't refresh access token; %s", err.Error()), true)
		}
	}

	return token.accessToken
}

// ForceRefreshToken should only be used on rare errors to immeaditalty refresh our access Token, regradles of our expirationTime
func (token *Token) ForceRefreshToken() error {
	return token.refreshAccessToken()
}

// refreshAccessToken uses the refreshToken to exchange for a new accessToken
func (token *Token) refreshAccessToken() error {
	parameters := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": token.refreshToken,
	}
	headers := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	responseMap, err := util.MakeHTTPRequest("POST", tokenURL, headers, parameters, nil)
	if err != nil {
		return fmt.Errorf("couldn't POST request refreshed token; %s", err.Error())
	}

	token.accessToken = responseMap["access_token"].(string)
	token.expirationTime = time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))-60) * time.Second)

	return nil
}
