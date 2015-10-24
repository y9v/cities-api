package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestCity(t *testing.T) {
	db := h.CreateDB(t)
	CreateCitiesBucket(db)
	CreateCountriesBucket(db)

	cityAttrs := []string{
		"Name",
		"DE",
		"10000000",
		"40.748817",
		"-73.958428",
		"Timezone",
	}
	cityString := strings.Join(cityAttrs, "\t")

	Convey("City from string", t, func() {
		Convey("When the string is correct", func() {
			city, err := cityFromString(1, cityString)

			Convey("Sets the city id from param", func() {
				So(city.ID, ShouldEqual, 1)
			})

			Convey("Sets the city attributes from the string", func() {
				So(city.Name, ShouldEqual, cityAttrs[0])
				So(city.CountryCode, ShouldEqual, cityAttrs[1])
				So(city.Population, ShouldEqual, 10000000)
				So(city.Latitude, ShouldEqual, 40.748817)
				So(city.Longitude, ShouldEqual, -73.958428)
				So(city.Timezone, ShouldEqual, cityAttrs[5])
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the string is incorrect", func() {
			city, err := cityFromString(1, "")

			Convey("Leaves the city id blank", func() {
				So(city.ID, ShouldEqual, 0)
			})

			Convey("Leaves the city attributes blank", func() {
				So(city.Name, ShouldEqual, "")
				So(city.CountryCode, ShouldEqual, "")
				So(city.Population, ShouldEqual, 0)
				So(city.Latitude, ShouldEqual, 0)
				So(city.Longitude, ShouldEqual, 0)
				So(city.Timezone, ShouldEqual, "")
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("City to string", t, func() {
		city := City{
			ID: 1, Name: "New York", CountryCode: "US", Population: 8600000,
			Latitude: 40.748817, Longitude: -73.985428, Timezone: "USA/New York",
		}

		Convey("Joins the city properties with tab chars", func() {
			expected := "New York\tUS\t8600000\t40.748817\t-73.985428\tUSA/New York"
			So(expected, ShouldEqual, city.toString())
		})
	})

	Convey("Find the city", t, func() {
		Convey("When the record exists", func() {
			h.PutToBucket(t, db, CitiesBucketName, "1", cityString)

			Convey("With includeCountry set to false", func() {
				city, err := FindCity(db, "1", false)

				Convey("Returns a city with attributes set", func() {
					expected, _ := cityFromString(1, cityString)
					So(city, ShouldResemble, expected)
				})

				Convey("Returns no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("With includeCountry set to true", func() {
				Convey("When country record exists", func() {
					countryAttrs := []string{"DE", "Germany", "en|Germany"}
					countryString := strings.Join(countryAttrs, "\t")
					h.PutToBucket(t, db, CountriesBucketName, "1", countryString)

					city, err := FindCity(db, "1", true)

					Convey("Returns a city with attributes set", func() {
						expected, _ := cityFromString(1, cityString)
						expected.Country, _ = countryFromString(1, countryString)
						So(city, ShouldResemble, expected)
					})

					Convey("Returns no error", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the record does not exist", func() {
			city, err := FindCity(db, "0", false)

			Convey("Returns a nil instead of a city", func() {
				So(city, ShouldBeNil)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When an incorrect value is stored in the db", func() {
			h.PutToBucket(t, db, CitiesBucketName, "2", "")
			city, err := FindCity(db, "2", false)

			Convey("Returns an empty city", func() {
				So(city, ShouldResemble, &City{})
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
