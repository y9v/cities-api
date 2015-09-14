package main

import (
	"fmt"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/parser"
	"log"
)

var options = config.Options{
	Port:               "8080",
	Timeout:            5,
	CORSOrigins:        []string{"http://localhost"},
	Locales:            []string{"en"},
	CitiesFile:         "data/cities.txt",
	AlternateNamesFile: "data/alternate.txt",
}

func main() {
	fmt.Printf("* Listening on port %s\n\n", options.Port)
	log.Fatal(Server().ListenAndServe())
}

func init() {
	fmt.Println("* Booting cities service...")
	fmt.Println("* Loading configuration...")
	config.Load(&options, "config.json")

	fmt.Println("* Connecting to the database...")
	InitDBSession()

	if GetAppStatus().IsIndexed() {
		fmt.Println("[PARSER] Skipping, already done")
		return
	} else {
		go parser.Scan(
			db, options.Locales, options.CitiesFile, options.AlternateNamesFile,
		)
	}
}
