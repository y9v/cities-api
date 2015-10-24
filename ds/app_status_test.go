package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppStatus(t *testing.T) {
	Convey("Test is indexed", t, func() {
		Convey("Returns true when status is ok", func() {
			appStatus := AppStatus{Statistics: &Statistics{Status: "ok"}}
			So(appStatus.IsIndexed(), ShouldBeTrue)
		})

		Convey("Returns false when status is not ok", func() {
			appStatus := AppStatus{}
			So(appStatus.IsIndexed(), ShouldBeFalse)
		})
	})

	Convey("Test get app status", t, func() {
		db := h.CreateDB(t)
		CreateStatisticsBucket(db)

		Convey("When indexing is done", func() {
			h.PutToBucket(t, db, StatisticsBucketName, "cities_count", "1000")
			h.PutToBucket(t, db, StatisticsBucketName, "city_names_count", "2000")

			appStatus := GetAppStatus(db)

			Convey("Sets app status to \"ok\"", func() {
				So(appStatus.Statistics.Status, ShouldEqual, "ok")
			})

			Convey("Has correct cities count", func() {
				So(appStatus.Statistics.CitiesCount, ShouldEqual, 1000)
			})

			Convey("Has correct city names count", func() {
				So(appStatus.Statistics.CityNamesCount, ShouldEqual, 2000)
			})
		})

		Convey("When still indexing", func() {
			h.DeleteFromBucket(t, db, StatisticsBucketName, "cities_count")
			h.DeleteFromBucket(t, db, StatisticsBucketName, "city_names_count")

			appStatus := GetAppStatus(db)

			Convey("Sets the app status to \"indexing\"", func() {
				So(appStatus.Statistics.Status, ShouldEqual, "indexing")
			})

			Convey("Has 0 for cities count", func() {
				So(appStatus.Statistics.CitiesCount, ShouldEqual, 0)
			})

			Convey("Has 0 for city names count", func() {
				So(appStatus.Statistics.CityNamesCount, ShouldEqual, 0)
			})
		})
	})
}
