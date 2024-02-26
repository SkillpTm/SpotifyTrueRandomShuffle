// Package api ...
package api

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const (
    responseType = "code"
    scopes = `user-read-playback-state%20user-read-currently-playing%20user-modify-playback-state%20user-read-private%20user-read-email%20playlist-read-private%20playlist-read-collaborative`
    authURL = "https://accounts.spotify.com/authorize?"
)

var (
    redirectURI = util.AppConfig.RedirectDomain + util.AppConfig.CallbackPort + util.AppConfig.CallbackPath

    state = util.GenerateRandomString(16)
    tokenChannel = make(chan *Token)
    UserToken *Token = nil
)

// <---------------------------------------------------------------------------------------------------->



// AuthUser authorizes our access to the user by getting the access token from Spotify
func AuthUser() {
    startHTTPServer()
    requestUserAuth()

    UserToken = <-tokenChannel
}



// startHTTPServer starts a server that listens and severs on a callback
func startHTTPServer() {
    http.HandleFunc(util.AppConfig.CallbackPath, handleAuthCode)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
    })
    go func() {
        err := http.ListenAndServe(util.AppConfig.CallbackPort, nil)
        if err != nil {
            util.LogError(err)
        }
    }()
}



// requestUserAuth prints the link required for the user auth
func requestUserAuth() {
    fmt.Printf("Please click this link and accept access: %sclient_id=%s&response_type=%s&redirect_uri=%s&state=%s&scope=%s&show_dialog=true", authURL, util.AppConfig.ClientID, responseType, redirectURI, state, scopes)
}



// handleAuthCode is a handler that handles the response from Spotify and puts a Token into the tokenChannel
func handleAuthCode(w http.ResponseWriter, r *http.Request) {

    query := r.URL.Query()
    callbackState := query.Get("state")

    if (state != callbackState) {
        http.NotFound(w, r)
        util.LogError(fmt.Errorf("state mismatch: %s != %s", state, callbackState))
        log.Fatalf("state mismatch: %s != %s\n", state, callbackState)
    }

    token, err := exchangeToken(query.Get("code"))
    if err != nil {
        http.Error(w, "couldn't get token", http.StatusForbidden)
        util.LogError(err)
        log.Fatal(err)
    }

    fmt.Fprintf(w, "Login Completed!")
    tokenChannel <- &token
}



// exchangeToken uses the authcode from Spotify to return the AccessToken, ExpirationTime and RefreshToken
func exchangeToken(authCode string) (Token, error) {
    parameters := map[string]string{
        "grant_type": "authorization_code",
        "code" : authCode,
        "redirect_uri" : redirectURI,
    }
    headers := map[string]string{
        "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
        "Content-Type" : "application/x-www-form-urlencoded",
    }

    responseMap, err := util.MakePOSTRequest("https://accounts.spotify.com/api/token", parameters, headers)
    if err != nil {
        return Token{}, err
    }

    return Token{
        AccessToken: responseMap["access_token"].(string),
        ExpirationTime: time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))) * time.Second),
        RefreshToken: responseMap["refresh_token"].(string),
    }, nil
}