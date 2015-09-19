package parser

import (
	"bufio"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	"os"
	"strconv"
	"strings"
)

func scanCities(db *bolt.DB, filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	citiesCount := 0

	err = db.Batch(func(tx *bolt.Tx) error {
		citiesBucket := tx.Bucket(ds.CitiesBucketName)
		cityNamesBucket := tx.Bucket(ds.CityNamesBucketName)

		for scanner.Scan() {
			cityData := strings.Split(scanner.Text(), "\t")
			cityBytes, err := prepareCityBytes(cityData)
			if err != nil {
				return err
			}

			citiesBucket.Put([]byte(cityData[0]), cityBytes)

			population, _ := strconv.ParseInt(cityData[14], 0, 64)

			addCityToIndex(
				cityNamesBucket, cityData[0], cityData[1], "", uint32(population),
			)

			citiesCount++
		}

		return err
	})

	return citiesCount, err
}
