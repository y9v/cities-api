package ds

import (
	"github.com/lebedev-yury/cities/config"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestSerializers(t *testing.T) {
	options := config.Options{Port: "9200"}

	Convey("City for serialization", t, func() {
		city := City{
			ID: 1, Name: "Venice", Population: 51298, CountryCode: "IT",
			Latitude: 45.43713, Longitude: 12.33265, Timezone: "Europe/Rome",
		}

		url := url.URL{Host: "", Scheme: ""}
		actual := city.ForSerialization(&url, &options)

		Convey("Has data attribute with a single item", func() {
			So(len(actual.Data), ShouldEqual, 1)
		})
	})

	Convey("Cities for serialization", t, func() {
		city := City{
			ID: 1, Name: "Venice", Population: 51298, CountryCode: "IT",
			Latitude: 45.43713, Longitude: 12.33265, Timezone: "Europe/Rome",
		}
		cities := Cities{Cities: []*City{&city}}

		url := url.URL{Host: "example.com", Scheme: "https"}
		actual := cities.ForSerialization(&url, &options)

		Convey("Has data attribute with a single item", func() {
			So(len(actual.Data), ShouldEqual, 1)
		})
	})
}
