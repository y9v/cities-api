package ds

import (
	"bytes"
	"github.com/boltdb/bolt"
	"sort"
)

type CityNames []*CityName

func appendCityName(slice CityNames, i *CityName) CityNames {
	for _, el := range slice {
		if el.CityId == i.CityId {
			return slice
		}
	}
	return append(slice, i)
}

func (cityNames *CityNames) Uniq() {
	var uniqCityNames CityNames
	for _, cityName := range *cityNames {
		uniqCityNames = appendCityName(uniqCityNames, cityName)
	}

	*cityNames = uniqCityNames
}

func (cityNames *CityNames) Limit(max int) {
	if len(*cityNames) > max {
		limitedCityNames := *cityNames
		*cityNames = limitedCityNames[:max]
	}
}

func searchCityNames(
	db *bolt.DB, locales []string, query string, limit int,
) (*CityNames, error) {
	var cityNames CityNames

	err := db.View(func(tx *bolt.Tx) error {
		var err error
		c := tx.Bucket(CityNamesBucketName).Cursor()

		prefix := []byte(PrepareCityNameKey(query))
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var cityName *CityName
			cityName, err = CityNameFromString(string(k), string(v))
			cityNames = append(cityNames, cityName)
		}

		return err
	})

	sort.Sort(CityNamesComparator{cityNames, locales})
	cityNames.Uniq()
	cityNames.Limit(limit)

	return &cityNames, err
}
