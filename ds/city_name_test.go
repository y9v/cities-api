package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

func TestCityName(t *testing.T) {
	Convey("Prepare cityname key", t, func() {
		Convey("Downcases the string", func() {
			So(PrepareCityNameKey("Foo"), ShouldEqual, "foo")
		})

		Convey("Removes whitespaces", func() {
			So(PrepareCityNameKey(" foo bar "), ShouldEqual, "foobar")
		})

		Convey("Removes dashes", func() {
			So(PrepareCityNameKey("foo-bar"), ShouldEqual, "foobar")
		})

		Convey("Removes pipes", func() {
			So(PrepareCityNameKey("foo|bar"), ShouldEqual, "foobar")
		})
	})

	Convey("Cityname to string", t, func() {
		cityName := CityName{
			Name: "New York", CityId: 1, Locale: "en", Population: 8600000,
		}

		Convey("Joins the cityname properties with tab chars", func() {
			So("New York\t1\ten\t8600000", ShouldEqual, cityName.toString())
		})
	})

	Convey("Cityname from string", t, func() {
		db := h.CreateDB(t)
		CreateCityNamesBucket(db)

		cityNameAttrs := []string{"Name", "1", "Locale", "10000000"}
		cityNameString := strings.Join(cityNameAttrs, "\t")

		Convey("When the string is correct", func() {
			cityName, err := CityNameFromString("key", cityNameString)

			Convey("Sets the key from param", func() {
				So(cityName.Key, ShouldEqual, "key")
			})

			Convey("Sets the cityname attributes from the string", func() {
				So(cityName.Name, ShouldEqual, cityNameAttrs[0])
				So(cityName.CityId, ShouldEqual, 1)
				So(cityName.Locale, ShouldEqual, cityNameAttrs[2])
			})

			Convey("Parses the population int from the string", func() {
				population, _ := strconv.ParseInt(cityNameAttrs[3], 0, 64)
				So(cityName.Population, ShouldEqual, uint32(population))
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the string is incorrect", func() {
			cityName, err := CityNameFromString("key", "")

			Convey("Returns an empty cityname", func() {
				So(cityName, ShouldResemble, &CityName{})
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
