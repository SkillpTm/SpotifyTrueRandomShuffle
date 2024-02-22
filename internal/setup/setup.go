// Package setup ...
package setup

// <---------------------------------------------------------------------------------------------------->

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// <---------------------------------------------------------------------------------------------------->

const (
	serverDomain = "http://localhost"
	serverPort = ":8080"
	serverCallback = "/callback"
)

const redirectURI = serverDomain + serverPort + serverCallback

var (
	_, _, _ = loadEnv()
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch = make(chan *spotify.Client)
	state = "abc123"
)

// <---------------------------------------------------------------------------------------------------->



func loadEnv() (string, string, string) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file", err)
		return "", "", ""
	}

	return os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"), os.Getenv("SPOTIFY_REDIRECT_URL")
}


func Setup() (*spotify.Client, *spotify.PrivateUser) {
	startHTTPServer()
	return loginUser()
}


func startHTTPServer() {
	http.HandleFunc(serverCallback, createClient)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(serverPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}


func loginUser() (*spotify.Client, *spotify.PrivateUser) {
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	return client, user
}


func createClient(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	fmt.Fprintf(w, "Login Completed!")
	ch <- spotify.New(auth.Client(r.Context(), tok))
}