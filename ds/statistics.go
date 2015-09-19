package ds

import (
	"github.com/boltdb/bolt"
	"strconv"
)

type Statistics struct {
	CountriesCount int `json:"countries_count"`
	CitiesCount    int `json:"cities_count"`
	CityNamesCount int `json:"city_names_count"`
}

func (statistics Statistics) Save(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(StatisticsBucketName)

		err := bucket.Put(
			[]byte("countries_count"),
			[]byte(strconv.Itoa(statistics.CountriesCount)),
		)

		err = bucket.Put(
			[]byte("cities_count"),
			[]byte(strconv.Itoa(statistics.CitiesCount)),
		)

		err = bucket.Put(
			[]byte("city_names_count"),
			[]byte(strconv.Itoa(statistics.CityNamesCount)),
		)

		return err
	})
}

func GetStatistics(db *bolt.DB) (*Statistics, error) {
	var stat Statistics

	err := db.View(func(tx *bolt.Tx) error {
		var err error

		if bucket := tx.Bucket(StatisticsBucketName); bucket != nil {
			stat.CountriesCount, err = strconv.Atoi(
				string(bucket.Get([]byte("countries_count"))),
			)

			stat.CitiesCount, err = strconv.Atoi(
				string(bucket.Get([]byte("cities_count"))),
			)

			stat.CityNamesCount, err = strconv.Atoi(
				string(bucket.Get([]byte("city_names_count"))),
			)
		}

		return err
	})

	return &stat, err
}
