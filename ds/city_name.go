package ds

import (
	"strconv"
	"strings"
)

type CityName struct {
	Key        string
	Name       string
	CityId     int
	Locale     string
	Population uint32
}

func (cityName *CityName) toString() string {
	return cityName.Name + "\t" + strconv.Itoa(cityName.CityId) + "\t" +
		cityName.Locale + "\t" + strconv.Itoa(int(cityName.Population))
}

func PrepareCityNameKey(key string) string {
	for _, s := range []string{" ", "|", "-"} {
		key = strings.Replace(key, s, "", -1)
	}

	return strings.ToLower(key)
}

func CityNameFromString(key string, cityNameString string) (*CityName, error) {
	var cityName CityName
	var err error

	cityNameData := strings.Split(cityNameString, "\t")

	if len(cityNameData) == 4 {
		var population int64
		population, err = strconv.ParseInt(cityNameData[3], 0, 64)
		cityName.CityId, err = strconv.Atoi(cityNameData[1])

		cityName.Key = key
		cityName.Name = cityNameData[0]
		cityName.Locale = cityNameData[2]
		cityName.Population = uint32(population)
	} else {
		err = InvalidDataError{CityNamesBucketName, key, cityNameString}
	}

	return &cityName, err
}
