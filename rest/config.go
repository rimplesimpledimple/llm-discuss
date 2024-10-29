package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port int `json:"port"`
}

func NewDefaultConfig() *config {
	return &config{
		Port: 8080,
	}
}

func loadConfig() config {
	// Define the file path for the configuration file
	configFilePath := "config.json"

	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// If the file does not exist, create it using the default configuration
		defaultConfig := NewDefaultConfig()
		configJSON, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling default config: %v", err)
		}

		// Write the default configuration to the file
		err = os.WriteFile(configFilePath, configJSON, 0644)
		if err != nil {
			log.Fatalf("Error writing default config to file: %v", err)
		}

		return *defaultConfig
	}

	// If the file exists, load the configuration from it
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var config config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return config
}
