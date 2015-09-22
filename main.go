package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/cache"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/ds"
	"github.com/lebedev-yury/cities/parser"
	"log"
)

func main() {
	fmt.Println("* Booting cities service...")

	fmt.Println("* Loading configuration...")
	options := config.Load("config.json")

	fmt.Println("* Connecting to the database...")
	db, err := bolt.Open("cities.db", 0600, nil)
	if err != nil {
		panic(fmt.Sprintf("[DB] Couldn't connect to the db: %v", err))
	}

	c := cache.New()

	if ds.GetAppStatus(db).IsIndexed() {
		fmt.Println("[PARSER] Skipping, already done")
	} else {
		go parser.Scan(
			db, options.Locales, options.MinPopulation,
			options.CountriesFile, options.CitiesFile,
			options.AlternateNamesFile,
		)
	}

	fmt.Printf("* Listening on port %s\n\n", options.Port)
	log.Fatal(Server(db, options, c).ListenAndServe())
}
