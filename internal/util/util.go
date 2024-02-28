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
	// open log file
	logFile, err := os.OpenFile(AppConfig.errorLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(errors.New("couldn't open log file: " + err.Error()))
	}
	defer logFile.Close()
	
	// write to log file
	fmt.Fprintf(logFile, "%v: %v\n", time.Now(), logErr)

	// exit the program
	log.Fatal(logErr)
}



// GenerateRandomString generates a random string of Base64 characters
func GenerateRandomString(length int) string {
	// create a byte slice
	bytes := make([]byte, length)

	// populate it with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		LogError(errors.New("couldn't rand read bytes for random string: " + err.Error()))
	}

	return base64.StdEncoding.EncodeToString(bytes)
}



// GetJSONData will provide a map with JSON data of the prvoided file
func GetJSONData(filePath string) (map[string]interface{}, error) {
	var jsonData map[string]interface{}

	// open JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't open JSON file (%s): %s", filePath, err.Error())
	}
	defer jsonFile.Close()

	// read JSON data from file
	rawJSONData, err := io.ReadAll(jsonFile)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't read JSON file (%s): %s", filePath, err.Error())
	}

	// convert JSON data to map
	err = json.Unmarshal(rawJSONData, &jsonData)
	if err != nil {
		return jsonData, errors.New("couldn't unmarshal raw JSON data: " + err.Error())
	}

	return jsonData, nil
}



// WriteJSONData will take a map with JSON data and the file path and write to that file
func WriteJSONData(filePath string, inputData map[string]interface{}) error {
	// open JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("couldn't open JSON file (%s): %s", filePath, err.Error())
	}
	defer jsonFile.Close()

    // marshal the map into JSON
    jsonData, err := json.MarshalIndent(inputData, "", "	")
    if err != nil {
        return errors.New("couldn't marshal JSON data: " + err.Error())
    }

    // write jsonData to file
    _, err = jsonFile.Write(jsonData)
    if err != nil {
        return fmt.Errorf("couldn't write JSON data to JSON file (%s): %s", filePath, err.Error())
    }

	return nil
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

	// execute the request
	response, err := httpClient.Do(request)
	if err != nil {
		return responseMap, errors.New("couldn't receive POST request response: " + err.Error())
	}
	defer response.Body.Close()

	// read the response
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

	// define responseMap for later use
	var responseMap map[string]interface{}

	// create inital request
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return responseMap, errors.New("couldn't create GET request: " + err.Error())
	}

	// add token header
	request.Header.Set("Authorization", "Bearer " + accessToken)

	// execute the request
	response, err := httpClient.Do(request)
	if err != nil {
		return responseMap, errors.New("couldn't receive GET request response: " + err.Error())
	}
	defer response.Body.Close()

	// read the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return responseMap, errors.New("couldn't read GET request response: " + err.Error())
	}

	// convert response to map, the response can be empty and still valid
	if (len(responseBody) > 0) {
		err = json.Unmarshal(responseBody, &responseMap)
		if err != nil {
			return responseMap, errors.New("couldn't unmarshal JSON GET request response body: " + err.Error())
		}
	}

	// check if we got an error code as a response
	_, notOK := responseMap["error"]
	if (notOK) {
		return responseMap, fmt.Errorf("spotify responded with an error %d to GET request: %s", int(responseMap["error"].(map[string]interface{})["status"].(float64)), responseMap["error"].(map[string]interface{})["message"].(string))
	}

	return responseMap, nil
}