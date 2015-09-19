package ds

import (
	"bytes"
	"github.com/boltdb/bolt"
	"strings"
)

type Country struct {
	Id           string            `json:"-"`
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	Translations map[string]string `json:"translations"`
}

func countryFromString(id string, countryString string) (*Country, error) {
	var country Country
	var err error

	countryData := strings.Split(countryString, "\t")

	if len(countryData) == 3 {
		country.Id = id
		country.Code = countryData[0]
		country.Name = countryData[1]

		country.Translations = make(map[string]string)
		for _, trData := range strings.Split(countryData[2], ";") {
			data := strings.Split(trData, "|")
			if len(data) == 2 {
				country.Translations[data[0]] = data[1]
			}
		}
	} else {
		err = InvalidDataError{CountriesBucketName, id, countryString}
	}

	return &country, err
}

func FindCountry(db *bolt.DB, id string) (*Country, error) {
	var country *Country = nil

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CountriesBucketName)
		val := bucket.Get([]byte(id))
		var err error

		if val != nil {
			country, err = countryFromString(id, string(val))
		}
		return err
	})

	return country, err
}

func FindCountryByCode(db *bolt.DB, code string) (*Country, error) {
	var country *Country = nil

	err := db.View(func(tx *bolt.Tx) error {
		var err error
		c := tx.Bucket(CountriesBucketName).Cursor()

		code := []byte(code)
		for k, v := c.First(); v != nil; k, v = c.Next() {
			if bytes.HasPrefix(v, code) {
				country, err = countryFromString(string(k), string(v))
				break
			}
		}

		return err
	})
	return country, err
}
