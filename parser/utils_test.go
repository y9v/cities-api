package parser

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestUtils(t *testing.T) {
	db := h.CreateDB(t)
	ds.CreateCityNamesBucket(db)
	ds.CreateCountriesBucket(db)

	Convey("Test prepare country bytes", t, func() {
		Convey("When the data is correct", func() {
			data := []string{
				"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"11", "12", "13", "14", "15", "16", "17", "18",
			}

			expected := []byte("0\t4\ten|4")
			actual, err := prepareCountryBytes(data)

			Convey("Joins the array in the right order with tabs", func() {
				So(actual, ShouldResemble, expected)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the data is incorrect", func() {
			data := []string{"yolo"}
			actual, err := prepareCountryBytes(data)

			Convey("Returns an empty bytes array", func() {
				var bytes []byte
				So(actual, ShouldResemble, bytes)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Test add translations to country", t, func() {
		translations := []string{"de|Deutschland", "ru|Германия"}

		countryAttrs := []string{"DE", "Germany", "en|Germany"}
		countryString := strings.Join(countryAttrs, "\t")
		h.PutToBucket(t, db, ds.CountriesBucketName, "1", countryString)

		err := db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket(ds.CountriesBucketName)
			return addTranslationsToCountry(b, 1, translations)
		})

		country, err := ds.FindCountry(db, "1")

		Convey("Does not modify country data", func() {
			So(country.Code, ShouldEqual, countryAttrs[0])
			So(country.Name, ShouldEqual, countryAttrs[1])
		})

		Convey("Keeps old translations", func() {
			So(country.Translations["en"], ShouldEqual, "Germany")
		})

		Convey("Adds new translations", func() {
			So(country.Translations["de"], ShouldEqual, "Deutschland")
			So(country.Translations["ru"], ShouldEqual, "Германия")
		})

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Test prepare city bytes", t, func() {
		Convey("When the data is correct", func() {
			data := []string{
				"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"11", "12", "13", "14", "15", "16", "17", "18",
			}

			expected := []byte("1\t8\t14\t4\t5\t17")
			actual, err := prepareCityBytes(data)

			Convey("Joins the array in the right order with tabs", func() {
				So(actual, ShouldResemble, expected)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the data is incorrect", func() {
			data := []string{"yolo"}
			actual, err := prepareCityBytes(data)

			Convey("Returns an empty bytes array", func() {
				var bytes []byte
				So(actual, ShouldResemble, bytes)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Test add city to index", t, func() {
		Convey("If no record for the key exist yet", func() {
			err := db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket(ds.CityNamesBucketName)
				return addCityToIndex(b, "1", "Berlin", "en", 2000000)
			})

			Convey("Puts the city name string to the bucket", func() {
				actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "berlin")
				So(actual, ShouldEqual, "Berlin\t1\ten\t2000000")
			})

			Convey("Not return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("If a record for the key exists", func() {
			existing := "Moscow\t1\ten\t12000000"
			h.PutToBucket(t, db, ds.CityNamesBucketName, "moscow", existing)

			err := db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket(ds.CityNamesBucketName)
				return addCityToIndex(b, "2", "Moscow", "en", 20000)
			})

			Convey("Does not overwrites the existing entry", func() {
				actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "moscow")
				So(actual, ShouldEqual, existing)
			})

			Convey("Adds city id as postfix for a new entry key", func() {
				actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "moscow|2")
				So(actual, ShouldEqual, "Moscow\t2\ten\t20000")
			})

			Convey("Not return an error", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Test is supported locale", t, func() {
		locales := []string{"ru", "en", "de"}

		Convey("Returns true for supported locales", func() {
			for _, locale := range locales {
				So(isSupportedLocale(locale, locales), ShouldBeTrue)
			}
		})

		Convey("Returns false for unsupported locales", func() {
			So(isSupportedLocale("jp", locales), ShouldBeFalse)
		})
	})
}
