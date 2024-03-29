// Package auth is responsible for handling everything to do with the authorization and token exchanges with Spotify
package auth

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/SkillpTm/SpotifyTrueRandomShuffle/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

const (
	contentType  = "application/x-www-form-urlencoded"
	responseType = "code"
	scopes       = "user-read-playback-state%20user-modify-playback-state%20playlist-read-private%20playlist-read-collaborative%20playlist-modify-public%20playlist-modify-private%20user-read-private%20user-read-email"

	authURL  = "https://accounts.spotify.com/authorize?"
	tokenURL = "https://accounts.spotify.com/api/token"
)

var (
	state               = util.GenerateRandomString(16)
	tokenChannel        = make(chan *Token)
	UserToken    *Token = nil
)

// <---------------------------------------------------------------------------------------------------->

// User authorizes our access to the user by getting the access token from Spotify
func User() {
	startHTTPServer()
	requestUserAuth()

	UserToken = <-tokenChannel
}

// startHTTPServer starts a server that listens and severs on a callback
func startHTTPServer() {
	// setup handlers
	http.HandleFunc(util.AppConfig.CallbackPath, handleAuthCode)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// user goroutine to listen and serve server
	go func() {
		err := http.ListenAndServe(util.AppConfig.CallbackPort, nil)
		if err != nil {
			util.LogError(fmt.Errorf("couldn't listen and server http server; %s", err.Error()), true)
		}
	}()
}

// requestUserAuth prints the link required for the user auth
func requestUserAuth() {
	fmt.Printf("Please click this link and accept access: \033[34m%sclient_id=%s&response_type=%s&redirect_uri=%s&state=%s&scope=%s&show_dialog=%t\033[0m", authURL, util.AppConfig.ClientID, responseType, util.AppConfig.RedirectURI, state, scopes, util.AppConfig.RequestAuthEveryTime)
}

// handleAuthCode is a handler that handles the response from Spotify and puts a Token into the tokenChannel
func handleAuthCode(w http.ResponseWriter, r *http.Request) {

	// query for state
	query := r.URL.Query()
	callbackState := query.Get("state")

	// compare the states
	if state != callbackState {
		http.Error(w, "state mismatch", http.StatusForbidden)
		util.LogError(fmt.Errorf("state mismatch: %s != %s", state, callbackState), true)
	}

	// exchange the token
	token, err := exchangeToken(query.Get("code"))
	if err != nil {
		http.Error(w, "couldn't get token", http.StatusForbidden)
		util.LogError(fmt.Errorf("couldn't exchange for access token; %s", err.Error()), true)
	}

	fmt.Fprintf(w, "Login Completed!")
	tokenChannel <- &token
}

// exchangeToken uses the authcode from Spotify to return the AccessToken, ExpirationTime and RefreshToken
func exchangeToken(authCode string) (Token, error) {
	parameters := map[string]string{
		"grant_type":   "authorization_code",
		"code":         authCode,
		"redirect_uri": util.AppConfig.RedirectURI,
	}
	headers := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(util.AppConfig.ClientID+":"+util.AppConfig.ClientSecret)),
		"Content-Type":  contentType,
	}

	// request token from Spotify
	responseMap, err := util.MakeHTTPRequest("POST", tokenURL, headers, parameters, nil)
	if err != nil {
		return Token{}, fmt.Errorf("couldn't POST request access token; %s", err.Error())
	}

	return Token{
		accessToken:    responseMap["access_token"].(string),
		expirationTime: time.Now().Add(time.Duration(int(responseMap["expires_in"].(float64))-60) * time.Second),
		refreshToken:   responseMap["refresh_token"].(string),
	}, nil
}
