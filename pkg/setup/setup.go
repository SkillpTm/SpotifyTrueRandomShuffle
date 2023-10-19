// Package setup ...
package setup

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)


func GetEnvs() (string, string, string) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file", err)
		return "", "", ""
	}

	redirectURL := os.Getenv("SPOTIFY_REDIRECT_URL")
	id := os.Getenv("SPOTIFY_ID")
	secret := os.Getenv("SPOTIFY_SECRET")

	return id, secret, redirectURL
}