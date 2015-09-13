package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strings"
)

func ParseAlternateNames() {
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
					err = addCityToIndex(
						cityNamesBucket, city.Id, cityData[3], cityData[2], city.Population,
					)
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

func addCityToIndex(
	bucket *bolt.Bucket, id string, name string, locale string, population string,
) error {
	cityNameKey := []byte(strings.ToLower(name))

	if confCityName := bucket.Get(cityNameKey); confCityName != nil {
		cityName := cityNameFromString(string(cityNameKey), string(confCityName))
		if cityName.CityId != id {
			cityNameKey = []byte(string(cityNameKey) + "|" + id)
		}
	}

	return bucket.Put(
		cityNameKey, []byte(name+"\t"+id+"\t"+locale+"\t"+population),
	)
}

func isSupportedLocale(locale string) bool {
	for _, item := range configuration.Locales {
		if item == locale {
			return true
		}
	}
	return false
}
