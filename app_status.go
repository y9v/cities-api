package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
)

type Stat struct {
	CitiesCount    int `json:"cities_count"`
	CityNamesCount int `json:"city_names_count"`
}

type AppStatus struct {
	Status string `json:"status"`
	Stats  Stat   `json:"statistics"`
}

var statsBucketName = []byte("statistics")

func (appStatus *AppStatus) IsOK() bool {
	return appStatus.Status == "ok"
}

func GetAppStatus() *AppStatus {
	var err error
	var appStatus AppStatus

	appStatus.Stats, err = getStats()
	if err != nil || appStatus.Stats.CitiesCount == 0 {
		appStatus.Status = "indexing"
	} else {
		appStatus.Status = "ok"
	}

	return &appStatus
}

func getStats() (Stat, error) {
	var stat Stat

	err := db.View(func(tx *bolt.Tx) error {
		var err error

		if bucket := tx.Bucket(statsBucketName); bucket != nil {
			stat.CitiesCount, err = strconv.Atoi(
				string(bucket.Get([]byte("cities_count"))),
			)

			stat.CityNamesCount, err = strconv.Atoi(
				string(bucket.Get([]byte("city_names_count"))),
			)
		}

		return err
	})

	return stat, err
}

func SaveStats(citiesCount int, cityNamesCount int) {
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(statsBucketName)
		bucket.Put([]byte("cities_count"), []byte(strconv.Itoa(citiesCount)))
		bucket.Put([]byte("city_names_count"), []byte(strconv.Itoa(cityNamesCount)))
		return nil
	})
}

func CreateStatsBucket() {
	db.Update(func(tx *bolt.Tx) error {
		fmt.Println("[DB] Creating bucket \"statistics\"...")
		tx.DeleteBucket(statsBucketName)
		_, err := tx.CreateBucket(statsBucketName)
		if err != nil {
			return fmt.Errorf("[DB] Error: %s", err)
		}
		return nil
	})
}
