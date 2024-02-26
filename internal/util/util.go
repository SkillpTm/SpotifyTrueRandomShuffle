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
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// <---------------------------------------------------------------------------------------------------->



func LogError(logErr error) {
    logFile, err := os.OpenFile("./logs/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()

    fmt.Fprintf(logFile, "%v: %v\n", time.Now(), logErr)
}


func GenerateRandomString(length int) string {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        LogError(err)
        log.Fatal(err)
    }

    return base64.StdEncoding.EncodeToString(bytes)
}


func MakePOSTRequest(requestURL string, parameters map[string]string, headers map[string]string) (map[string]interface{}, error) {
    httpClient := &http.Client{}
    postBody := url.Values{}

    for key, value := range parameters {
        postBody.Set(key, value)
    }

    request, err := http.NewRequest("POST", requestURL, strings.NewReader(postBody.Encode()))
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't create POST request: " + err.Error())
    }

    for key, value := range headers {
        request.Header.Set(key, value)
    }

    response, err := httpClient.Do(request)
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't receive POST request response: " + err.Error())
    }
    defer response.Body.Close()

    responseBody, err := io.ReadAll(response.Body)
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't read POST request response: " + err.Error())
    }

    var responseMap map[string]interface{}

    if (len(responseBody) > 0) {
        err = json.Unmarshal(responseBody, &responseMap)
        if err != nil {
            return map[string]interface{}{}, errors.New("couldn't unmarshal JSON POST request response body: " + err.Error())
        }
    }

    return responseMap, nil
}

func MakeGETRequest(requestURL string, accessToken string) (map[string]interface{}, error) {
    httpClient := &http.Client{}

    request, err := http.NewRequest("GET", requestURL, nil)
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't create GET request: " + err.Error())
    }

    request.Header.Set("Authorization", "Bearer " + accessToken)

    response, err := httpClient.Do(request)
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't receive GET request response: " + err.Error())
    }
    defer response.Body.Close()

    responseBody, err := io.ReadAll(response.Body)
    if err != nil {
        return map[string]interface{}{}, errors.New("couldn't read GET request response: " + err.Error())
    }

    var responseMap map[string]interface{}

    if (len(responseBody) > 0) {
        err = json.Unmarshal(responseBody, &responseMap)
        if err != nil {
            return map[string]interface{}{}, errors.New("couldn't unmarshal JSON GET request response body: " + err.Error())
        }
    }

    return responseMap, nil
}