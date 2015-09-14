package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestStatistics(t *testing.T) {
	db := h.CreateDB(t)
	CreateStatisticsBucket(db)

	Convey("Save statistics", t, func() {
		statistics := Statistics{CitiesCount: 100000, CityNamesCount: 200000}
		err := statistics.Save(db)

		Convey("Saves the cities count to the db", func() {
			val, _ := strconv.Atoi(
				h.ReadFromBucket(t, db, StatisticsBucketName, "cities_count"),
			)

			So(val, ShouldEqual, statistics.CitiesCount)
		})

		Convey("Saves the city names count to the db", func() {
			val, _ := strconv.Atoi(
				h.ReadFromBucket(t, db, StatisticsBucketName, "city_names_count"),
			)

			So(val, ShouldEqual, statistics.CityNamesCount)
		})

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Get statistics", t, func() {
		Convey("When the values are in the db", func() {
			h.PutToBucket(t, db, StatisticsBucketName, "cities_count", "5000")
			h.PutToBucket(t, db, StatisticsBucketName, "city_names_count", "9000")

			statistics := GetStatistics(db)

			Convey("Reads the cities count from db", func() {
				So(statistics.CitiesCount, ShouldEqual, 5000)
			})

			Convey("Reads the city names count from db", func() {
				So(statistics.CityNamesCount, ShouldEqual, 9000)
			})
		})

		Convey("When the values are not in the db", func() {
			h.DeleteFromBucket(t, db, StatisticsBucketName, "cities_count")
			h.DeleteFromBucket(t, db, StatisticsBucketName, "city_names_count")

			Convey("Does not panic", func() {
				So(func() { GetStatistics(db) }, ShouldNotPanic)
			})
		})
	})
}
