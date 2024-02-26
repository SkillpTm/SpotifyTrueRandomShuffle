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

// importConfig imports to config file from ./configs/config
func importConfig() (map[string]interface{}, error) {
    var configData map[string]interface{}

    // import config file
    configFile, err := os.Open("./configs/config.json")
    if err != nil {
        return configData, errors.New("couldn't open config file: " + err.Error())
    }
    defer configFile.Close()

    // read config data from file
    rawConfigData, err := io.ReadAll(configFile)
    if err != nil {
        return configData, errors.New("couldn't read config file: " + err.Error())
    }

    // convert data to map
    err = json.Unmarshal(rawConfigData, &configData)
    if err != nil {
        return configData, errors.New("couldn't unmarshal raw config data: " + err.Error())
    }

    return configData, nil
}