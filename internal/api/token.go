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
    accessToken string
    expirationTime time.Time
    refreshToken string
}

// <---------------------------------------------------------------------------------------------------->



func (token *Token) GetAccessToken() string {
    currentTime := time.Now()

    if (token.expirationTime.Before(currentTime)) {
        err := token.refreshAccessToken(currentTime)
        if err != nil {
            util.LogError(err)
            log.Fatal(err)
        }
    }

    return token.accessToken
}


func (token *Token) refreshAccessToken(currentTime time.Time) error {
    parameters := map[string]string{
        "grant_type": "refresh_token",
        "refresh_token" : token.refreshToken,
    }
    headers := map[string]string{
        "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
        "Content-Type" : "application/x-www-form-urlencoded",
    }

    responseMap, err := util.MakePOSTRequest("https://accounts.spotify.com/api/token", parameters, headers)
    if err != nil {
        return err
    }

    token.accessToken = responseMap["access_token"].(string)
    token.expirationTime = time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))) * time.Second)

    return nil
}