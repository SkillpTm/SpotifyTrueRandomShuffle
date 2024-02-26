// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/joho/godotenv"
)

// <---------------------------------------------------------------------------------------------------->

var AppConfig Config

// <---------------------------------------------------------------------------------------------------->



// Config struct is a type to hold our config data
type Config struct {
	CallbackPath string
	CallbackPort string
	LoopRefreshTime float64
	envPath string
	errorLogPath string

	ClientID string
	ClientSecret string
	RedirectDomain string
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



// importConfig loads ./configs/config onto AppConfig
func importConfig() error {

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
		CallbackPath: configData["callbackPath"].(string),
		CallbackPort: configData["callbackPort"].(string),
		LoopRefreshTime: configData["loopRefreshTime"].(float64),
		envPath: configData["paths"].(map[string]string)["env"],
		errorLogPath: configData["paths"].(map[string]string)["errorLog"],
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

    return nil
}