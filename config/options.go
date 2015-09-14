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

func Load(options *Options, filename string) {
	if envFilename := os.Getenv("CONFIG"); len(envFilename) > 0 {
		filename = envFilename
	}

	if file, err := os.Open(filename); err == nil {
		err = json.NewDecoder(file).Decode(&options)
		if err != nil {
			fmt.Println("[CONFIG] Error parsing configuration:", err)
			os.Exit(1)
		}
	}

	if envPort := os.Getenv("PORT"); len(envPort) > 0 {
		options.Port = envPort
	}
}
