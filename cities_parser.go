package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strings"
	"time"
)

func ParseCities() {
	file, err := os.Open(configuration.CitiesFile)
	if err != nil {
		fmt.Println("[PARSER] Error opening cities file:", err)
		os.Exit(1)
	}
	defer file.Close()

	if GetAppStatus().IsOK() {
		fmt.Println("[PARSER] Skipping, already done")
		return
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	CreateCitiesBucket()
	CreateCityNamesBucket()
	CreateStatsBucket()

	startTime := time.Now()
	fmt.Println("[PARSER] Started cities parsing")

	citiesCount := 0
	err = db.Batch(func(tx *bolt.Tx) error {
		citiesBucket := tx.Bucket(citiesBucketName)
		cityNamesBucket := tx.Bucket(cityNamesBucketName)

		for scanner.Scan() {
			cityData := strings.Split(scanner.Text(), "\t")
			id := []byte(cityData[0])

			err = citiesBucket.Put(id, prepareCityBytes(cityData))
			if err != nil {
				return err
			}

			err = addCityToIndex(
				cityNamesBucket, cityData[0], cityData[1], "", cityData[14],
			)
			if err != nil {
				return err
			}

			citiesCount++
		}

		return err
	})

	fmt.Printf("[PARSER] Added %d cities\n", citiesCount)
	cityNamesCount := ParseAlternateNames()
	fmt.Printf("[PARSER] Added %d city names\n", citiesCount+cityNamesCount)
	SaveStats(citiesCount, citiesCount+cityNamesCount)

	if err != nil {
		fmt.Println("[PARSER] Error:", err)
	} else {
		fmt.Printf("[PARSER] Parsing done (in %s)\n", time.Since(startTime))
	}
}

func prepareCityBytes(cityData []string) []byte {
	result := cityData[1] + "\t" + cityData[8] + "\t" + cityData[14] +
		"\t" + cityData[4] + "\t" + cityData[5] + "\t" + cityData[17]

	return []byte(result)
}