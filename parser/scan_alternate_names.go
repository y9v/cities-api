package parser

import (
	"bufio"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	"os"
	"strconv"
	"strings"
)

func scanAlternateNames(
	db *bolt.DB, filename string, locales []string,
) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	cityNamesCount := 0
	countriesTranslations := make(map[int]map[string]string)

	err = db.Update(func(tx *bolt.Tx) error {
		cityNamesBucket := tx.Bucket(ds.CityNamesBucketName)
		countriesBucket := tx.Bucket(ds.CountriesBucketName)

		for scanner.Scan() {
			nameData := strings.Split(scanner.Text(), "\t")

			if isSupportedLocale(nameData[2], locales) || nameData[4] == "1" {
				city, _ := ds.FindCity(db, nameData[1], false)
				if city != nil {
					addCityToIndex(
						cityNamesBucket, strconv.Itoa(city.ID), nameData[3],
						nameData[2], city.Population,
					)

					cityNamesCount++
				} else if nameData[2] != "en" {
					country, _ := ds.FindCountry(db, nameData[1])
					if country != nil {
						if countriesTranslations[country.ID] == nil {
							countriesTranslations[country.ID] = make(map[string]string)
						}
						if countriesTranslations[country.ID][nameData[2]] == "" {
							countriesTranslations[country.ID][nameData[2]] = nameData[3]
						}
					}
				}
			}
		}

		for id, translations := range countriesTranslations {
			var values []string
			for locale, name := range translations {
				values = append(values, locale+"|"+name)
			}
			addTranslationsToCountry(countriesBucket, id, values)
		}

		return err
	})

	return cityNamesCount, err
}
