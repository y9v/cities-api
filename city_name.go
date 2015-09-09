package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"sort"
	"strconv"
	"strings"
)

type CityName struct {
	Key        string
	Name       string
	CityId     string
	Population uint32
}

type CityNames []CityName

func (slice CityNames) Len() int {
	return len(slice)
}
func (slice CityNames) Less(i, j int) bool {
	return slice[i].Population > slice[j].Population
}
func (slice CityNames) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var cityNamesBucketName = []byte("city_names")

func (cityNames *CityNames) Limit(max int) {
	if len(*cityNames) > max {
		limitedCityNames := *cityNames
		*cityNames = limitedCityNames[:max]
	}
}

func cityNameFromString(key string, cityNameString string) *CityName {
	cityNameData := strings.Split(cityNameString, "\t")
	population, _ := strconv.ParseInt(cityNameData[2], 0, 64)

	return &CityName{
		Key:        key,
		Name:       cityNameData[0],
		CityId:     cityNameData[1],
		Population: uint32(population),
	}
}

func CreateCityNamesSchema() {
	db.Update(func(tx *bolt.Tx) error {
		fmt.Println("* [DB] Creating bucket \"city_names\"...")
		_, err := tx.CreateBucket([]byte(cityNamesBucketName))
		if err != nil {
			return fmt.Errorf("* [DB] Error: %s", err)
		}
		return nil
	})
}

func SearchCityNames(query string, limit int) (*CityNames, error) {
	var cityNames CityNames

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(cityNamesBucketName).Cursor()

		prefix := []byte(query)
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			cityName := cityNameFromString(string(k), string(v))
			cityNames = append(cityNames, *cityName)
		}

		return nil
	})

	sort.Sort(cityNames)
	cityNames.Limit(limit)

	return &cityNames, err
}
