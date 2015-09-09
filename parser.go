package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strings"
	"time"
)

func ParseCitiesFile() {
	file, err := os.Open(configuration.CitiesFile)
	if err != nil {
		fmt.Println("* [PARSER] Error opening cities file:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	CreateCitiesSchema()
	CreateCityNamesSchema()

	startTime := time.Now()
	fmt.Println("* [PARSER] Started cities parsing")

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

			addCityToIndex(cityNamesBucket, cityData[0], cityData[1], cityData[14])
			if err != nil {
				return err
			}
		}

		return err
	})

	ParseAlternateNamesFile()

	if err != nil {
		fmt.Println("* [PARSER] Error:", err)
	} else {
		fmt.Printf("* [PARSER] Parsing done (in %s)\n", time.Since(startTime))
	}
}

func ParseAlternateNamesFile() {
	file, err := os.Open(configuration.AlternateNamesFile)
	if err != nil {
		fmt.Println("* [PARSER] Error opening alternate names file:", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	fmt.Println("* [PARSER] Started alternate names parsing")

	err = db.Batch(func(tx *bolt.Tx) error {
		cityNamesBucket := tx.Bucket(cityNamesBucketName)

		for scanner.Scan() {
			cityData := strings.Split(scanner.Text(), "\t")

			if isSupportedLocale(cityData[2]) || cityData[4] == "1" {
				city, err := FindCity(cityData[1])
				if err != nil {
					return err
				}

				if city != nil {
					err = addCityToIndex(cityNamesBucket, city.Id, cityData[3], city.Population)
					if err != nil {
						return err
					}
				}
			}
		}

		return err
	})

	if err != nil {
		fmt.Println("* [PARSER] Error:", err)
	}
}

func prepareCityBytes(cityData []string) []byte {
	result := cityData[1] + "\t" + cityData[8] + "\t" + cityData[14] +
		"\t" + cityData[4] + "\t" + cityData[5] + "\t" + cityData[17]

	return []byte(result)
}

func addCityToIndex(bucket *bolt.Bucket, id string, name string, population string) error {
	cityNameKey := []byte(strings.ToLower(name))

	if confCityName := bucket.Get(cityNameKey); confCityName != nil {
		cityName := cityNameFromString(string(cityNameKey), string(confCityName))
		if cityName.CityId != id {
			cityNameKey = []byte(string(cityNameKey) + "|" + id)
		}
	}

	return bucket.Put(cityNameKey, []byte(name+"\t"+id+"\t"+population))
}

func isSupportedLocale(locale string) bool {
	for _, item := range configuration.Locales {
		if item == locale {
			return true
		}
	}
	return false
}
