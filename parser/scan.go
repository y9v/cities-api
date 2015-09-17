package parser

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	"time"
)

func Scan(
	db *bolt.DB, locales []string, citiesFile string, alternateNamesFile string,
) {
	startTime := time.Now()

	ds.CreateCitiesBucket(db)
	ds.CreateCityNamesBucket(db)
	ds.CreateStatisticsBucket(db)

	fmt.Println("[PARSER] Started cities parsing")
	var cityNamesCount int

	citiesCount, err := scanCities(db, citiesFile)
	if err == nil {
		fmt.Println("[PARSER] Started alternate names parsing")
		cityNamesCount, err = scanAlternateNames(db, alternateNamesFile, locales)
	}

	if err != nil {
		panic(fmt.Sprintf("[PARSER] Error: %v", err))
	} else {
		ds.Statistics{
			CitiesCount:    citiesCount,
			CityNamesCount: citiesCount + cityNamesCount,
		}.Save(db)

		fmt.Printf("[PARSER] Added %d cities\n", citiesCount)
		fmt.Printf("[PARSER] Added %d city names\n", citiesCount+cityNamesCount)
		fmt.Printf("[PARSER] Parsing done (in %s)\n", time.Since(startTime))
	}
}
