package parser

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
)

func prepareCityBytes(cityData []string) []byte {
	return []byte(
		cityData[1] + "\t" + cityData[8] + "\t" + cityData[14] +
			"\t" + cityData[4] + "\t" + cityData[5] + "\t" + cityData[17],
	)
}

func addCityToIndex(
	bucket *bolt.Bucket, id string, name string, locale string, population string,
) error {
	var err error
	var cityName *ds.CityName

	cityNameKey := []byte(ds.PrepareCityNameKey(name))
	if conflict := bucket.Get(cityNameKey); conflict != nil {
		cityName, err = ds.CityNameFromString(string(cityNameKey), string(conflict))
		if cityName.CityId != id {
			cityNameKey = []byte(string(cityNameKey) + "|" + id)
		}
	}

	err = bucket.Put(
		cityNameKey, []byte(name+"\t"+id+"\t"+locale+"\t"+population),
	)

	return err
}

func isSupportedLocale(locale string, locales []string) bool {
	for _, item := range locales {
		if item == locale {
			return true
		}
	}
	return false
}
