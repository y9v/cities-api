package parser

import (
	"github.com/lebedev-yury/cities/ds"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScan(t *testing.T) {
	Convey("Test scan", t, func() {
		db := h.CreateDB(t)
		locales := []string{"ru", "en"}

		countriesFilename := h.CreateTempfile(t, "WS\tWSM\t882\tWS\tSamoa\tApia\t2944\t192001\tOC\t.ws\tWST\tTala\t685\t\t\tsm,en-WS\t4034894\t\t")
		citiesFilename := h.CreateTempfile(t, "890516\tGwanda\tGwanda\tJawunda\t-20.93333\t29\tP\tPPLA\tZW\t\t07\t\t\t\t14450\t\t982\tAfrica/Harare\t2009-06-30")
		alternateNamesFilename := h.CreateTempfile(t, "10\t890516\tru\tГуанда\t\t\t\t")

		Convey("When both files are present and valid", func() {
			done := make(chan bool, 1)
			Scan(db, done, locales, 2000, countriesFilename, citiesFilename, alternateNamesFilename)

			Convey("Does not panics", func() {
				done := make(chan bool, 1)

				So(func() {
					Scan(db, done, locales, 2000, countriesFilename, citiesFilename, alternateNamesFilename)
				}, ShouldNotPanic)
			})

			Convey("Writes countries count to the statistics bucket", func() {
				actual := h.ReadFromBucket(t, db, ds.StatisticsBucketName, "countries_count")
				So(actual, ShouldEqual, "1")
			})

			Convey("Writes cities count to the statistics bucket", func() {
				actual := h.ReadFromBucket(t, db, ds.StatisticsBucketName, "cities_count")
				So(actual, ShouldEqual, "1")
			})

			Convey("Writes citynames count to the statistics bucket", func() {
				actual := h.ReadFromBucket(t, db, ds.StatisticsBucketName, "city_names_count")
				So(actual, ShouldEqual, "2")
			})
		})

		Convey("When the countries file does not exist", func() {
			Convey("Panics", func() {
				done := make(chan bool, 1)

				So(func() {
					Scan(db, done, locales, 2000, "fake.txt", citiesFilename, alternateNamesFilename)
				}, ShouldPanic)
			})
		})

		Convey("When the cities file does not exist", func() {
			Convey("Panics", func() {
				done := make(chan bool, 1)

				So(func() {
					Scan(db, done, locales, 2000, countriesFilename, "fake.txt", alternateNamesFilename)
				}, ShouldPanic)
			})
		})

		Convey("When the alternate names file does not exist", func() {
			Convey("Panics", func() {
				done := make(chan bool, 1)

				So(func() {
					Scan(db, done, locales, 2000, countriesFilename, citiesFilename, "fake.txt")
				}, ShouldPanic)
			})
		})
	})
}
