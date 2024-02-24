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
    serverPort = ":8080"
    serverCallback = "/callback"

    responseType = "code"
    scopes = `user-read-playback-state%20user-read-currently-playing%20user-modify-playback-state%20user-read-private%20user-read-email%20playlist-read-private%20playlist-read-collaborative`
    authURL = "https://accounts.spotify.com/authorize?"
)

var (
    clientID, clientSecret, serverDomain = util.LoadEnv()
    redirectURI = serverDomain + serverPort + serverCallback

    state = util.GenerateRandomString(16)
    tokenChannel = make(chan *Token)
    UserToken *Token = nil
)

// <---------------------------------------------------------------------------------------------------->



func AuthUser() {
    startHTTPServer()
    requestUserAuth()

    UserToken = <-tokenChannel
}


func startHTTPServer() {
    http.HandleFunc(serverCallback, handleAuthCode)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.NotFound(w, r)
    })
    go func() {
        err := http.ListenAndServe(serverPort, nil)
        if err != nil {
            util.LogError(err)
            log.Fatal(err)
        }
    }()
}


func requestUserAuth() {
    fmt.Println("Please click this link and accept access: " + fmt.Sprintf("%sclient_id=%s&response_type=%s&redirect_uri=%s&state=%s&scope=%s&show_dialog=true", authURL, clientID, responseType, redirectURI, state, scopes))
}


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


func exchangeToken(authCode string) (Token, error) {
    parameters := map[string]string{
        "grant_type": "authorization_code",
        "code" : authCode,
        "redirect_uri" : redirectURI,
    }
    headers := map[string]string{
        "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)),
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