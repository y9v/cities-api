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
	parsingDone := make(chan bool, 1)

	if ds.GetAppStatus(db).IsIndexed() {
		fmt.Println("[PARSER] Skipping, already done")
		parsingDone <- true
	} else {
		go parser.Scan(
			db, parsingDone, options.Locales, options.MinPopulation,
			options.CountriesFile, options.CitiesFile,
			options.AlternateNamesFile,
		)
	}

	<-parsingDone
	fmt.Println("[CACHE] Warming up...")
	warmUpSearchCache(db, c, options.Locales, 5)
	fmt.Println("[CACHE] Warming up done")

	fmt.Printf("* Listening on port %s\n\n", options.Port)
	log.Fatal(Server(db, options, c).ListenAndServe())
}
