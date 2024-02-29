// Package util ...
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
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
	envPath string
	errorLogPath string
	LoopRefreshTime float64
	RequestAuthEveryTime bool
	ShufflePlaylistPath string
	ShufflePlaylistSize int

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



// importConfig loads ./configs/config onto AppConfig
func importConfig() error {

	configData, err := GetJSONData("./configs/config.json")
	if err != nil {
		return fmt.Errorf("couldn't get config.json; %s", err.Error())
	}

	// set configData to exportable var
	AppConfig = Config{
		CallbackPath: configData["callbackPath"].(string),
		CallbackPort: configData["callbackPort"].(string),
		envPath: configData["paths"].(map[string]interface{})["env"].(string),
		errorLogPath: configData["paths"].(map[string]interface{})["errorLog"].(string),
		LoopRefreshTime: configData["loopRefreshTime"].(float64),
		RequestAuthEveryTime: configData["requestAuthEveryTime"].(bool),
		ShufflePlaylistPath: configData["paths"].(map[string]interface{})["shufflePlaylist"].(string),
		ShufflePlaylistSize : int(configData["shufflePlaylistSize"].(float64)),
	}

	return nil
}



// loadEnv imports the envs for the spotify API from the .env file on AppConfig
func loadEnv() error {
	// load envs into enviroment
	err := godotenv.Load(AppConfig.envPath)
	if err != nil {
		return fmt.Errorf("couldn't load .env file; %s", err.Error())
	}

	AppConfig.ClientID = os.Getenv("SPOTIFY_ID")
	AppConfig.ClientSecret = os.Getenv("SPOTIFY_SECRET")
	AppConfig.RedirectDomain = os.Getenv("SPOTIFY_REDIRECT_DOMAIN")
	AppConfig.RedirectURI = AppConfig.RedirectDomain + AppConfig.CallbackPort + AppConfig.CallbackPath

	return nil
}