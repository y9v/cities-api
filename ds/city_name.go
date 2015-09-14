package ds

import (
	"bytes"
	"github.com/boltdb/bolt"
	"sort"
	"strconv"
	"strings"
)

type CityName struct {
	Key        string
	Name       string
	CityId     string
	Locale     string
	Population uint32
}

type CityNames []*CityName

func (cityNames *CityNames) Limit(max int) {
	if len(*cityNames) > max {
		limitedCityNames := *cityNames
		*cityNames = limitedCityNames[:max]
	}
}

func appendCityName(slice CityNames, i *CityName) CityNames {
	for _, el := range slice {
		if el.CityId == i.CityId {
			return slice
		}
	}
	return append(slice, i)
}

func (cityNames CityNames) Uniq() {
	var uniqCityNames CityNames
	for _, cityName := range cityNames {
		uniqCityNames = appendCityName(uniqCityNames, cityName)
	}

	cityNames = uniqCityNames
}

func PrepareCityNameKey(key string) string {
	for _, s := range []string{" ", "|", "-"} {
		key = strings.Replace(key, s, "", -1)
	}

	return strings.ToLower(key)
}

func CityNameFromString(key string, cityNameString string) *CityName {
	cityNameData := strings.Split(cityNameString, "\t")
	population, _ := strconv.ParseInt(cityNameData[3], 0, 64)

	return &CityName{
		Key:        key,
		Name:       cityNameData[0],
		CityId:     cityNameData[1],
		Locale:     cityNameData[2],
		Population: uint32(population),
	}
}

func SearchCityNames(
	db *bolt.DB, locales []string, query string, limit int,
) (*CityNames, error) {
	var cityNames CityNames

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(CityNamesBucketName).Cursor()

		prefix := []byte(PrepareCityNameKey(query))
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			cityName := CityNameFromString(string(k), string(v))
			cityNames = append(cityNames, cityName)
		}

		return nil
	})

	sort.Sort(CityNameComparator{cityNames, locales})
	cityNames.Uniq()
	cityNames.Limit(limit)

	return &cityNames, err
}
