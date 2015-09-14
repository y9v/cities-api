package ds

import (
	"github.com/boltdb/bolt"
)

var CitiesBucketName = []byte("cities")
var CityNamesBucketName = []byte("city_names")
var StatisticsBucketName = []byte("statistics")

func CreateCitiesBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(CitiesBucketName)
		_, err := tx.CreateBucket(CitiesBucketName)
		return err
	})
}

func CreateCityNamesBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(CityNamesBucketName)
		_, err := tx.CreateBucket(CityNamesBucketName)
		return err
	})
}

func CreateStatisticsBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(StatisticsBucketName)
		_, err := tx.CreateBucket(StatisticsBucketName)
		return err
	})
}
