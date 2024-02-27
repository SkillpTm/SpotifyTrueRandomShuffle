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

// LogError writes any error to a log file and then uses log.Fatal
func LogError(logErr error) {
    logFile, err := os.OpenFile(AppConfig.errorLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(errors.New("couldn't open log file: " + err.Error()))
    }
    defer logFile.Close()

    fmt.Fprintf(logFile, "%v: %v\n", time.Now(), logErr)

    log.Fatal(logErr)
}



// GenerateRandomString generates a random string of Base64 characters
func GenerateRandomString(length int) string {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        LogError(errors.New("couldn't rand read bytes for random string: " + err.Error()))
    }

    return base64.StdEncoding.EncodeToString(bytes)
}



// MakePOSTRequest makes a POST request with headers, parameters or plain text body data and returns the JSON format response as a map
func MakePOSTRequest(requestURL string, headers map[string]string, parameters map[string]string, bodyData map[string]interface{}) (map[string]interface{}, error) {
    httpClient := &http.Client{}

    // define vars for later use
    var responseMap map[string]interface{}
    var requestBody io.Reader

    // add parameters to requestBody
    if (len(parameters) != 0) {
        requestParameters := url.Values{}

        for key, value := range parameters {
            requestParameters.Set(key, value)
        }

        requestBody = strings.NewReader(requestParameters.Encode())

    // add bodyData to requestBody
    } else if (len(bodyData) != 0) {
        jsonBodyData, err := json.Marshal(bodyData)
        if err != nil {
            return responseMap, errors.New("couldn't marshal body data for POST request: " + err.Error())
        }

        requestBody = strings.NewReader(string(jsonBodyData))
    }

    // create request with request body
    request, err := http.NewRequest("POST", requestURL, requestBody)
    if err != nil {
        return responseMap, errors.New("couldn't create POST request: " + err.Error())
    }

    // add headers to request
    for key, value := range headers {
        request.Header.Set(key, value)
    }

    // execute request
    response, err := httpClient.Do(request)
    if err != nil {
        return responseMap, errors.New("couldn't receive POST request response: " + err.Error())
    }
    defer response.Body.Close()

    // read response
    responseBody, err := io.ReadAll(response.Body)
    if err != nil {
        return responseMap, errors.New("couldn't read POST request response: " + err.Error())
    }

    // convert response to map, the response can be empty and still valid
    if (len(responseBody) > 0) {
        err = json.Unmarshal(responseBody, &responseMap)
        if err != nil {
            return responseMap, errors.New("couldn't unmarshal JSON POST request response body: " + err.Error())
        }
    }

    // check if we got an error code as a response
    _, notOK := responseMap["error"]
    if (notOK) {
        return responseMap, fmt.Errorf("spotify responded with an error %d to POST request: %s", int(responseMap["error"].(map[string]interface{})["status"].(float64)), responseMap["error"].(map[string]interface{})["message"].(string))
    }

    return responseMap, nil
}



// MakeGETRequest makes a GET request with headers and returns the JSON format response as a map
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

    _, notOK := responseMap["error"]
    if (notOK) {
        return responseMap, fmt.Errorf("spotify responded with an error %d to GET request: %s", int(responseMap["error"].(map[string]interface{})["status"].(float64)), responseMap["error"].(map[string]interface{})["message"].(string))
    }

    return responseMap, nil
}