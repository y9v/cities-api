package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/ds"
	"github.com/lebedev-yury/cities/parser"
	"log"
)

var (
	db *bolt.DB

	options = config.Options{
		Port:               "8080",
		Timeout:            5,
		CORSOrigins:        []string{"http://localhost"},
		Locales:            []string{"en"},
		CitiesFile:         "data/cities.txt",
		AlternateNamesFile: "data/alternate.txt",
	}
)

func main() {
	fmt.Println("* Booting cities service...")
	fmt.Println("* Loading configuration...")
	config.Load(&options, "config.json")

	fmt.Println("* Connecting to the database...")
	InitDBSession()

	if ds.GetAppStatus(db).IsIndexed() {
		fmt.Println("[PARSER] Skipping, already done")
		return
	} else {
		go parser.Scan(
			db, options.Locales, options.CitiesFile, options.AlternateNamesFile,
		)
	}

	fmt.Printf("* Listening on port %s\n\n", options.Port)
	log.Fatal(Server().ListenAndServe())
}
