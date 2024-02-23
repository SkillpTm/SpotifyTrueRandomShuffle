// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// <---------------------------------------------------------------------------------------------------->



func LoadEnv() (string, string, string) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file", err)
		return "", "", ""
	}

	return os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"), os.Getenv("SPOTIFY_REDIRECT_URL")
}


func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
       panic(err)
    }

    return base64.StdEncoding.EncodeToString(bytes)
}


func MakePostRequest(requestURL string, parameters map[string]string, headers map[string]string) (map[string]interface{}, error) {
    httpClient := &http.Client{}
	postBody := url.Values{}

    for key, value := range parameters {
        postBody.Set(key, value)
    }

    request, err := http.NewRequest("POST", requestURL, strings.NewReader(postBody.Encode()))
    if err != nil {
        return map[string]interface{}{}, errors.New("Couldn't create POST request: " + err.Error())
    }

    for key, value := range headers {
        request.Header.Set(key, value)
    }

    response, err := httpClient.Do(request)
    if err != nil {
        return map[string]interface{}{}, errors.New("Couldn't request POST request: " + err.Error())
    }
    defer response.Body.Close()

    responseBody, err := io.ReadAll(response.Body)
    if err != nil {
        return map[string]interface{}{}, errors.New("Couldn't read POST request response: " + err.Error())
    }

    var responseMap map[string]interface{}

    err = json.Unmarshal(responseBody, &responseMap)
    if err != nil {
        return map[string]interface{}{}, errors.New("Couldn't unmarshal JSON POST request response: " + err.Error())
    }

	return responseMap, nil
}