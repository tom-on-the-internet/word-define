package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func getConfig() (Config, error) {
	var config Config

	err := createConfigIfNotExists()
	if err != nil {
		return config, err
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}

	appConfigDir := userConfigDir + "/word-define"
	configFilename := appConfigDir + "/config.json"

	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		// json was invalid, write empty json
		configJSON, _ := json.MarshalIndent(config, "", "  ")
		err = ioutil.WriteFile(configFilename, configJSON, 0600)
	}

	return config, err
}

func createConfigIfNotExists() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	// create directory if not exists
	_ = os.Mkdir(userConfigDir, 0700)

	appConfigdir := userConfigDir + "/word-define"
	configFilename := appConfigdir + "/config.json"

	// create directory if not exists
	_ = os.Mkdir(appConfigdir, 0700)

	// create file if not exists
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		var config Config

		configJSON, _ := json.MarshalIndent(config, "", "  ")

		return ioutil.WriteFile(configFilename, configJSON, 0600)
	}

	return nil
}
