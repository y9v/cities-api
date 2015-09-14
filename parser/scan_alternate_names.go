package parser

import (
	"bufio"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	"os"
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

	err = db.Batch(func(tx *bolt.Tx) error {
		cityNamesBucket := tx.Bucket(ds.CityNamesBucketName)

		for scanner.Scan() {
			cityData := strings.Split(scanner.Text(), "\t")

			if isSupportedLocale(cityData[2], locales) || cityData[4] == "1" {
				city, err := ds.FindCity(db, cityData[1])
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

					cityNamesCount++
				}
			}
		}

		return err
	})

	return cityNamesCount, err
}
