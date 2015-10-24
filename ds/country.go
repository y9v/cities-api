package ds

import (
	"bytes"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

type Country struct {
	ID           int               `json:"-"`
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	Translations map[string]string `json:"translations"`
}

func countryFromString(id int, countryString string) (*Country, error) {
	var country Country
	var err error

	countryData := strings.Split(countryString, "\t")

	if len(countryData) == 3 {
		country.ID = id
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
		err = InvalidDataError{CountriesBucketName, strconv.Itoa(id), countryString}
	}

	return &country, err
}

func FindCountry(db *bolt.DB, key string) (*Country, error) {
	var country *Country = nil

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CountriesBucketName)
		val := bucket.Get([]byte(key))
		var id int
		var err error

		if val != nil {
			if id, err = strconv.Atoi(key); err == nil {
				country, err = countryFromString(id, string(val))
			}
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
				var id int
				if id, err = strconv.Atoi(string(k)); err == nil {
					country, err = countryFromString(id, string(v))
				}
				break
			}
		}

		return err
	})
	return country, err
}
