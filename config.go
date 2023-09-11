package main

import (
	"encoding/json"
	"os"
)

func readConfig() error {
	configFile, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer configFile.Close()

	var config struct {
		Token     string   `json:"token"`
		Favorites []string `json:"fav"`
		Action    string   `json:"action"`
	}

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	authToken = config.Token
	registrationHeaders["Authorization"] = authToken

	if len(config.Favorites) > 0 {
		favoriteCourses = config.Favorites
	}

	if config.Action != "" {
		config.Action = "add"
	}

	action = config.Action

	return nil
}