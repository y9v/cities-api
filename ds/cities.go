package ds

import (
	"github.com/boltdb/bolt"
)

type Cities struct {
	Cities []*City `json:"cities,omitempty"`
}

func appendCity(db *bolt.DB, cities []*City, city *City, locale string) []*City {
	for _, el := range cities {
		if el.Name == city.Name {
			if country, _ := FindCountryByCode(db, city.CountryCode); country != nil {
				city.Name = city.Name + ", " + country.Translations[locale]
				break
			} else {
				return cities
			}
		}
	}
	return append(cities, city)
}

func SearchCities(
	db *bolt.DB, locales []string, query string, limit int,
) (*Cities, error) {
	var cities Cities

	cityNames, err := searchCityNames(db, locales, query, limit)

	var city *City
	for _, cityName := range *cityNames {
		city, err = FindCity(db, cityName.CityId, false)
		city.Name = cityName.Name
		cities.Cities = appendCity(db, cities.Cities, city, cityName.Locale)
	}

	return &cities, err
}
