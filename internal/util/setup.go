// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// <---------------------------------------------------------------------------------------------------->

var AppConfig Config

// <---------------------------------------------------------------------------------------------------->



// Config struct is a type to hold our config data
type Config struct {
	callbackPath string
	callbackPort string
	loopRefreshTime float64
	envPath string
	errorLogPath string
}



// importConfig imports to config file from ./configs/config
func importConfig() (error) {

    // import config file
    configFile, err := os.Open("./configs/config.json")
    if err != nil {
        return errors.New("couldn't open config file: " + err.Error())
    }
    defer configFile.Close()

    // read config data from file
    rawConfigData, err := io.ReadAll(configFile)
    if err != nil {
        return errors.New("couldn't read config file: " + err.Error())
    }

	var configData map[string]interface{}

    // convert data to map
    err = json.Unmarshal(rawConfigData, &configData)
    if err != nil {
        return errors.New("couldn't unmarshal raw config data: " + err.Error())
    }

	// set configData to exportable var
	AppConfig = Config{
		callbackPath: configData["callbackPath"].(string),
		callbackPort: configData["callbackPort"].(string),
		loopRefreshTime: configData["loopRefreshTime"].(float64),
		envPath: configData["paths"].(map[string]string)["env"],
		errorLogPath: configData["paths"].(map[string]string)["errorLog"],
	}

    return nil
}