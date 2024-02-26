// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

// <---------------------------------------------------------------------------------------------------->

var AppConfig Config

// <---------------------------------------------------------------------------------------------------->



// Config is a type to hold our config data
type Config struct {
	CallbackPath string
	CallbackPort string
	LoopRefreshTime float64
	envPath string
	errorLogPath string
	TempPlaylistPath string

	ClientID string
	ClientSecret string
	RedirectDomain string

	RedirectURI string
}



// Setup loads our config and envs onto AppConfig
func Setup() error {
	err := importConfig()
    if err != nil {
        return err
    }

	err = loadEnv()
    if err != nil {
        return err
    }

	return nil
}



// GetJSONData will provide a map with JSON data of the prvoided file
func GetJSONData(filePath string) (map[string]interface{}, error) {
	var configData map[string]interface{}

    // open JSON file
    configFile, err := os.Open(filePath)
    if err != nil {
        return configData, fmt.Errorf("couldn't open JSON file (%s): %s", filePath, err.Error())
    }
    defer configFile.Close()

    // read JSON data from file
    rawConfigData, err := io.ReadAll(configFile)
    if err != nil {
        return configData, fmt.Errorf("couldn't read JSON file (%s): %s", filePath, err.Error())
    }

    // convert JSON data to map
    err = json.Unmarshal(rawConfigData, &configData)
    if err != nil {
        return configData, errors.New("couldn't unmarshal raw JSON data: " + err.Error())
    }

	return configData, nil
}



// importConfig loads ./configs/config onto AppConfig
func importConfig() error {

	configData, err := GetJSONData("./configs/config.json")
    if err != nil {
        return errors.New("couldn't get config.json: " + err.Error())
    }

	// set configData to exportable var
	AppConfig = Config{
		CallbackPath: configData["callbackPath"].(string),
		CallbackPort: configData["callbackPort"].(string),
		LoopRefreshTime: configData["loopRefreshTime"].(float64),
		envPath: configData["paths"].(map[string]string)["env"],
		errorLogPath: configData["paths"].(map[string]string)["errorLog"],
		TempPlaylistPath: configData["paths"].(map[string]string)["tempPlaylist"],
	}

    return nil
}



// loadEnv imports the envs for the spotify API from the .env file on AppConfig
func loadEnv() error {
	// load envs into enviroment
    err := godotenv.Load(AppConfig.envPath)
    if err != nil {
        return errors.New("couldn't load .env file: " + err.Error())
    }

	AppConfig.ClientID = os.Getenv("SPOTIFY_ID")
	AppConfig.ClientSecret = os.Getenv("SPOTIFY_SECRET")
	AppConfig.RedirectDomain = os.Getenv("SPOTIFY_REDIRECT_DOMAIN")
	AppConfig.RedirectURI = AppConfig.RedirectDomain + AppConfig.CallbackPort + AppConfig.CallbackPath

    return nil
}