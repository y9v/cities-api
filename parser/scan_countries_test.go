package parser

import (
	"github.com/lebedev-yury/cities/ds"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanCountries(t *testing.T) {
	Convey("Test scan countries", t, func() {
		db := h.CreateDB(t)
		ds.CreateCountriesBucket(db)

		Convey("When cities files exists", func() {
			filename := h.CreateTempfile(t, "# Comment\nWS\tWSM\t882\tWS\tSamoa\tApia\t2944\t192001\tOC\t.ws\tWST\tTala\t685\t\t\tsm,en-WS\t4034894\t\t")

			count, err := scanCountries(db, filename)

			Convey("Stores parsed countries to the db", func() {
				actual := h.ReadFromBucket(t, db, ds.CountriesBucketName, "4034894")
				So(actual, ShouldEqual, "WS\tSamoa\ten|Samoa")
			})

			Convey("Returns number of scanned records", func() {
				So(count, ShouldEqual, 1)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the file has invalid data", func() {
			filename := h.CreateTempfile(t, "crap\ncrap\ncrap")
			count, err := scanCountries(db, filename)

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When some locale has blank id", func() {
			filename := h.CreateTempfile(t, "WS\tWSM\t882\tWS\tSamoa\tApia\t2944\t192001\tOC\t.ws\tWST\tTala\t685\t\t\tsm,en-WS\t\t\t")
			count, _ := scanCountries(db, filename)

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})
		})

		Convey("When countries file does not exist", func() {
			count, err := scanCountries(db, "fake.txt")

			Convey("Returns a zero number of scanned records", func() {
				So(count, ShouldEqual, 0)
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
