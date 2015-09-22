package cache

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCache(t *testing.T) {
	Convey("New", t, func() {
		actual := New()

		Convey("Returns a pointer to a cache instance", func() {
			So(actual, ShouldHaveSameTypeAs, &Cache{})
		})
	})

	Convey("Set and Get", t, func() {
		cache := New()
		expected := "value"
		key := "key"
		cache.Set(key, expected)

		Convey("Find a stored object by key", func() {
			actual, ok := cache.Get(key)
			So(actual, ShouldEqual, expected)
			So(ok, ShouldBeTrue)
		})

		Convey("Finds nothing if nothing was stored", func() {
			actual, ok := cache.Get("yolo")
			So(actual, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})
	})
}
