package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Options struct {
	Port               string
	Timeout            int
	CORSOrigins        []string
	Locales            []string
	CitiesFile         string
	AlternateNamesFile string
}

func buildDefault() *Options {
	return &Options{
		Port:               "8080",
		Timeout:            5,
		CORSOrigins:        []string{"http://localhost"},
		Locales:            []string{"en"},
		CitiesFile:         "data/cities.txt",
		AlternateNamesFile: "data/alternateNames.txt",
	}
}

func Load(filename string) *Options {
	options := buildDefault()

	if envFilename := os.Getenv("CONFIG"); len(envFilename) > 0 {
		filename = envFilename
	}

	if file, err := os.Open(filename); err == nil {
		err = json.NewDecoder(file).Decode(options)
		if err != nil {
			panic(fmt.Sprintf("Error parsing configuration: %v", err))
		}
	}

	if envPort := os.Getenv("PORT"); len(envPort) > 0 {
		options.Port = envPort
	}

	return options
}
