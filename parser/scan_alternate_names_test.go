package parser

import (
	"github.com/lebedev-yury/cities/ds"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanAlternateNames(t *testing.T) {
	Convey("Test scan alternate names", t, func() {
		db := h.CreateDB(t)
		ds.CreateCitiesBucket(db)
		ds.CreateCityNamesBucket(db)

		h.PutToBucket(t, db, ds.CitiesBucketName, "1", "Montreal\t\t\t\t\t")
		h.PutToBucket(t, db, ds.CitiesBucketName, "2", "Moscow\t\t\t\t\t")

		locales := []string{"de", "ru"}

		Convey("When alternate names file exists", func() {
			filename := h.CreateTempfile(
				t,
				"10\t1\tfr\tMontréal\t\t\t\t\n11\t2\tde\tMoskau\t\t\t\t\n12\t2\tru\tМосква\t\t\t\t13\t3\tde\tMünchen\t\t\t\t",
			)

			count, err := scanAlternateNames(db, filename, locales)

			Convey("Returns number of scanned records", func() {
				So(count, ShouldEqual, 2)
			})

			Convey("When the locale is supported", func() {
				Convey("Stores the record if the city exists", func() {
					actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "moskau")
					So(actual, ShouldEqual, "Moskau\t2\tde\t")
				})

				Convey("Doesn't store the record if the city doesn't exist", func() {
					actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "münchen")
					So(actual, ShouldEqual, "")
				})
			})

			Convey("When the locale is not supported", func() {
				Convey("Doesn't store the record", func() {
					actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "montréal")
					So(actual, ShouldEqual, "")
				})
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When alternate names file does not exist", func() {
			count, err := scanAlternateNames(db, "fake.txt", locales)

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
