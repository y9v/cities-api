package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Configuration struct {
	Port               string
	Timeout            int
	CORSOrigins        []string
	CORSHeader         string
	Locales            []string
	CitiesFile         string
	AlternateNamesFile string
}

var configuration = Configuration{
	Port:    "8080",
	Timeout: 2,
}

func LoadConfiguration(configuration *Configuration) {
	fileName := "config.json"
	if envFileName := os.Getenv("CONFIG"); len(envFileName) > 0 {
		fileName = envFileName
	}

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}

	err = json.NewDecoder(file).Decode(&configuration)

	if err != nil {
		fmt.Println("Error parsing config file: ", err)
		os.Exit(1)
	}

	if envPort := os.Getenv("PORT"); len(envPort) > 0 {
		configuration.Port = envPort
	}

	configuration.CORSHeader = strings.Join(configuration.CORSOrigins, ",")
}
