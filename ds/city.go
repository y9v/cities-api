package ds

import (
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

type City struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	CountryCode string   `json:"-"`
	Population  uint32   `json:"population"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Timezone    string   `json:"timezone"`
	Country     *Country `json:"country,omitempty"`
}

func cityFromString(id int, cityString string) (*City, error) {
	var city City
	var err error

	cityData := strings.Split(cityString, "\t")

	if len(cityData) == 6 {
		var population int64
		population, err = strconv.ParseInt(cityData[2], 0, 64)
		city.Latitude, err = strconv.ParseFloat(cityData[3], 64)
		city.Longitude, err = strconv.ParseFloat(cityData[4], 64)

		city.ID = id
		city.Name = cityData[0]
		city.CountryCode = cityData[1]
		city.Population = uint32(population)
		city.Timezone = cityData[5]
	} else {
		err = InvalidDataError{CitiesBucketName, strconv.Itoa(id), cityString}
	}

	return &city, err
}

func (city *City) toString() string {
	return city.Name + "\t" + city.CountryCode + "\t" +
		strconv.Itoa(int(city.Population)) + "\t" +
		strconv.FormatFloat(city.Latitude, 'f', 6, 64) + "\t" +
		strconv.FormatFloat(city.Longitude, 'f', 6, 64) + "\t" +
		city.Timezone
}

func FindCity(db *bolt.DB, key string, includeCountry bool) (*City, error) {
	var city *City = nil

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CitiesBucketName)
		val := bucket.Get([]byte(key))
		var err error

		if val != nil {
			id, _ := strconv.Atoi(key)
			city, err = cityFromString(id, string(val))
			if err == nil && includeCountry == true {
				city.Country, err = FindCountryByCode(db, city.CountryCode)
			}
		}
		return err
	})

	return city, err
}
