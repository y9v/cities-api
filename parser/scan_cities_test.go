package parser

import (
	"github.com/lebedev-yury/cities/ds"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanCities(t *testing.T) {
	Convey("Test scan cities", t, func() {
		db := h.CreateDB(t)
		ds.CreateCitiesBucket(db)
		ds.CreateCityNamesBucket(db)

		Convey("When cities files exists", func() {
			filename := h.CreateTempfile(t, "890516\tGwanda\tGwanda\tJawunda\t-20.93333\t29\tP\tPPLA\tZW\t\t07\t\t\t\t14450\t\t982\tAfrica/Harare\t2009-06-30\n890983\tGokwe\tGokwe\tGokwe\t-18.20476\t28.9349\tP\tPPL\tZW\t\t02\t\t\t\t18942\t\t1237\tAfrica/Harare\t2012-05-05")

			count, err := scanCities(db, filename)

			Convey("Stores parsed cities to the db", func() {
				actual := h.ReadFromBucket(t, db, ds.CitiesBucketName, "890516")
				So(actual, ShouldEqual, "Gwanda\tZW\t14450\t-20.93333\t29\tAfrica/Harare")

				actual = h.ReadFromBucket(t, db, ds.CitiesBucketName, "890983")
				So(actual, ShouldEqual, "Gokwe\tZW\t18942\t-18.20476\t28.9349\tAfrica/Harare")
			})

			Convey("Stores parsed city names to the db", func() {
				actual := h.ReadFromBucket(t, db, ds.CityNamesBucketName, "gwanda")
				So(actual, ShouldEqual, "Gwanda\t890516\ten\t14450")

				actual = h.ReadFromBucket(t, db, ds.CityNamesBucketName, "gokwe")
				So(actual, ShouldEqual, "Gokwe\t890983\ten\t18942")
			})

			Convey("Returns number of scanned records", func() {
				So(count, ShouldEqual, 2)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the file has invalid data", func() {
			filename := h.CreateTempfile(t, "crap\ncrap\ncrap")
			count, err := scanCities(db, filename)

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When cities file does not exist", func() {
			count, err := scanCities(db, "fake.txt")

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
