package ds

import (
	"github.com/boltdb/bolt"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBuckets(t *testing.T) {
	Convey("Cities bucket creation", t, func() {
		Convey("Creates a bucket", func() {
			db := h.CreateDB(t)
			CreateCitiesBucket(db)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(CitiesBucketName)
				So(b, ShouldNotBeNil)
				return nil
			})
		})

		Convey("Deletes the existing bucket", func() {
			db := h.CreateDB(t)
			oldBucket := h.CreateBucket(t, db, CitiesBucketName)
			CreateCitiesBucket(db)

			db.View(func(tx *bolt.Tx) error {
				newBucket := tx.Bucket(CitiesBucketName)
				So(oldBucket, ShouldNotEqual, newBucket)
				return nil
			})
		})
	})

	Convey("City names bucket creation", t, func() {
		Convey("Creates a bucket", func() {
			db := h.CreateDB(t)
			CreateCityNamesBucket(db)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(CityNamesBucketName)
				So(b, ShouldNotBeNil)
				return nil
			})
		})

		Convey("Deletes the existing bucket", func() {
			db := h.CreateDB(t)
			oldBucket := h.CreateBucket(t, db, CityNamesBucketName)
			CreateCitiesBucket(db)

			db.View(func(tx *bolt.Tx) error {
				newBucket := tx.Bucket(CityNamesBucketName)
				So(oldBucket, ShouldNotEqual, newBucket)
				return nil
			})
		})
	})

	Convey("Statistics bucket creation", t, func() {
		Convey("Creates a bucket", func() {
			db := h.CreateDB(t)
			CreateStatisticsBucket(db)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(StatisticsBucketName)
				So(b, ShouldNotBeNil)
				return nil
			})
		})

		Convey("Deletes the existing bucket", func() {
			db := h.CreateDB(t)
			oldBucket := h.CreateBucket(t, db, StatisticsBucketName)
			CreateCitiesBucket(db)

			db.View(func(tx *bolt.Tx) error {
				newBucket := tx.Bucket(StatisticsBucketName)
				So(oldBucket, ShouldNotEqual, newBucket)
				return nil
			})
		})
	})
}
