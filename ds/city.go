package ds

import (
	"github.com/boltdb/bolt"
	"strings"
)

type City struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Population  string `json:"population"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Timezone    string `json:"timezone"`
}

type Cities struct {
	Cities []*City `json:"cities,omitempty"`
}

func cityFromString(id string, cityString string) (*City, error) {
	var city City
	var err error

	cityData := strings.Split(cityString, "\t")

	if len(cityData) == 6 {
		city.Id = id
		city.Name = cityData[0]
		city.CountryCode = cityData[1]
		city.Population = cityData[2]
		city.Latitude = cityData[3]
		city.Longitude = cityData[4]
		city.Timezone = cityData[5]
	} else {
		err = InvalidDataError{CitiesBucketName, id, cityString}
	}

	return &city, err
}

func (city *City) toString() string {
	return city.Name + "\t" + city.CountryCode + "\t" + city.Population + "\t" +
		city.Latitude + "\t" + city.Longitude + "\t" + city.Timezone
}

func FindCity(db *bolt.DB, id string) (*City, error) {
	var city *City = nil

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CitiesBucketName)
		val := bucket.Get([]byte(id))
		var err error

		if val != nil {
			city, err = cityFromString(id, string(val))
		}
		return err
	})

	return city, err
}

func SearchCities(
	db *bolt.DB, locales []string, query string, limit int,
) (*Cities, error) {
	var cities Cities

	cityNames, err := searchCityNames(db, locales, query, limit)

	var city *City
	for _, cityName := range *cityNames {
		city, err = FindCity(db, cityName.CityId)
		city.Name = cityName.Name
		cities.Cities = append(cities.Cities, city)
	}

	return &cities, err
}
