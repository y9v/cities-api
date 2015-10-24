package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestCountry(t *testing.T) {
	db := h.CreateDB(t)
	CreateCountriesBucket(db)

	countryAttrs := []string{"DE", "Germany", "en|Germany;de|Deutschland"}
	countryString := strings.Join(countryAttrs, "\t")

	Convey("Country from string", t, func() {
		Convey("When the string is correct", func() {
			country, err := countryFromString(1, countryString)

			Convey("Sets the country id from param", func() {
				So(country.ID, ShouldEqual, 1)
			})

			Convey("Sets the country attributes", func() {
				So(country.Code, ShouldEqual, countryAttrs[0])
				So(country.Name, ShouldEqual, countryAttrs[1])
			})

			Convey("Sets the country translations", func() {
				So(country.Translations["en"], ShouldEqual, "Germany")
				So(country.Translations["de"], ShouldEqual, "Deutschland")
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the string is incorrect", func() {
			country, err := countryFromString(1, "")

			Convey("Leaves the country id blank", func() {
				So(country.ID, ShouldEqual, 0)
			})

			Convey("Leaves the country attributes blank", func() {
				So(country.Code, ShouldEqual, "")
				So(country.Name, ShouldEqual, "")
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Find the country", t, func() {
		Convey("When the record exists", func() {
			h.PutToBucket(t, db, CountriesBucketName, "1", countryString)
			country, err := FindCountry(db, "1")

			Convey("Returns a country with attributes set", func() {
				expected, _ := countryFromString(1, countryString)
				So(country, ShouldResemble, expected)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the record does not exist", func() {
			country, err := FindCountry(db, "0")

			Convey("Returns a nil instead of a country", func() {
				So(country, ShouldBeNil)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When an incorrect value is stored in the db", func() {
			h.PutToBucket(t, db, CountriesBucketName, "2", "")
			country, err := FindCountry(db, "2")

			Convey("Returns an empty country", func() {
				So(country, ShouldResemble, &Country{})
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Find the country by code", t, func() {
		Convey("When the record exists", func() {
			h.PutToBucket(t, db, CountriesBucketName, "1", countryString)
			country, err := FindCountryByCode(db, countryAttrs[0])

			Convey("Returns a country with attributes set", func() {
				expected, _ := countryFromString(1, countryString)
				So(country, ShouldResemble, expected)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the record does not exist", func() {
			country, err := FindCountryByCode(db, "HU")

			Convey("Returns a nil instead of a country", func() {
				So(country, ShouldBeNil)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When an incorrect value is stored in the db", func() {
			h.PutToBucket(t, db, CountriesBucketName, "2", "ES")
			country, err := FindCountryByCode(db, "ES")

			Convey("Returns an empty country", func() {
				So(country, ShouldResemble, &Country{})
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
