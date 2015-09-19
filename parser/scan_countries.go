package parser

import (
	"bufio"
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	"os"
	"strings"
)

func scanCountries(db *bolt.DB, filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	countriesCount := 0

	err = db.Batch(func(tx *bolt.Tx) error {
		countriesBucket := tx.Bucket(ds.CountriesBucketName)

		for scanner.Scan() {
			countryString := scanner.Text()
			if strings.HasPrefix(countryString, "#") {
				continue
			}

			countryData := strings.Split(countryString, "\t")
			countryBytes, err := prepareCountryBytes(countryData)
			if err != nil {
				return err
			}

			if id := countryData[16]; id != "" {
				countriesBucket.Put([]byte(id), countryBytes)
				countriesCount++
			}
		}

		return err
	})

	return countriesCount, err
}
