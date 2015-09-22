package ds

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/cache"
	"unicode/utf8"
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

func CachedCitiesSearch(
	db *bolt.DB, c *cache.Cache, locales []string, query string, limit int,
) (*Cities, error) {
	if utf8.RuneCountInString(query) > 2 {
		return SearchCities(db, locales, query, limit)
	}

	var err error
	var cities *Cities

	cacheKey := "cs." + string(limit) + "." + query
	if i, ok := c.Get(cacheKey); ok {
		cities = i.(*Cities)
	} else {
		cities, err = SearchCities(db, locales, query, limit)
		if len(cities.Cities) > 0 {
			c.Set(cacheKey, cities)
		}
	}

	return cities, err
}
