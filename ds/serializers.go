package ds

import (
	"github.com/lebedev-yury/cities/config"
	"net/url"
	"strconv"
)

type DataSerializer struct {
	Data []interface{} `json:"data"`
}

type CitySerializer struct {
	ID    int               `json:"id"`
	Type  string            `json:"type"`
	City  *City             `json:"attributes"`
	Links map[string]string `json:"links"`
}

func baseURL(url *url.URL, options *config.Options) string {
	if url.Host != "" {
		return url.Scheme + "://" + url.Host
	} else {
		return "http://localhost:" + options.Port
	}
}

func (city *City) url(url *url.URL, options *config.Options) string {
	return baseURL(url, options) + "/1.0/cities/" + strconv.Itoa(city.ID)
}

func (city *City) serializer(
	url *url.URL, options *config.Options,
) *CitySerializer {
	return &CitySerializer{
		ID:   city.ID,
		Type: "cities",
		City: city,
		Links: map[string]string{
			"self": city.url(url, options),
		},
	}
}

func (city *City) ForSerialization(
	url *url.URL, options *config.Options,
) *DataSerializer {
	return &DataSerializer{
		Data: []interface{}{city.serializer(url, options)},
	}
}

func (cities *Cities) ForSerialization(
	url *url.URL, options *config.Options,
) *DataSerializer {
	var citiesForSerialization = []interface{}{}
	for _, city := range cities.Cities {
		citiesForSerialization = append(
			citiesForSerialization, city.serializer(url, options),
		)
	}

	return &DataSerializer{Data: citiesForSerialization}
}
