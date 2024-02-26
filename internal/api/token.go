// Package api ...
package api

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"log"
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
        err := token.refreshAccessToken(currentTime)
        if err != nil {
            util.LogError(err)
            log.Fatal(err)
        }
    }

    return token.AccessToken
}


func (token *Token) refreshAccessToken(currentTime time.Time) error {
    parameters := map[string]string{
        "grant_type": "refresh_token",
        "refresh_token" : token.RefreshToken,
    }
    headers := map[string]string{
        "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
        "Content-Type" : "application/x-www-form-urlencoded",
    }

    responseMap, err := util.MakePOSTRequest("https://accounts.spotify.com/api/token", parameters, headers)
    if err != nil {
        return err
    }

    token.AccessToken = responseMap["access_token"].(string)
    token.ExpirationTime = time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))) * time.Second)

    return nil
}