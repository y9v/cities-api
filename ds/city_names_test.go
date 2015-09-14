package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCityNames(t *testing.T) {
	Convey("Limit citynames", t, func() {
		Convey("Limits a collection if too big", func() {
			actual := CityNames{
				&CityName{Name: "A"}, &CityName{Name: "B"}, &CityName{Name: "C"},
			}

			expected := CityNames{&CityName{Name: "A"}, &CityName{Name: "B"}}

			actual.Limit(2)
			So(actual, ShouldResemble, expected)
		})

		Convey("Does not changes the collection if not too big", func() {
			actual := CityNames{&CityName{Name: "A"}, &CityName{Name: "B"}}
			expected := actual
			actual.Limit(2)
			So(actual, ShouldResemble, expected)
		})
	})

	Convey("Uniq citynames", t, func() {
		Convey("Removes values duplicated by city id", func() {
			actual := CityNames{
				&CityName{Name: "Moscow", CityId: "1"},
				&CityName{Name: "Moskau", CityId: "1"},
				&CityName{Name: "Montreal", CityId: "2"},
			}

			expected := CityNames{
				&CityName{Name: "Moscow", CityId: "1"},
				&CityName{Name: "Montreal", CityId: "2"},
			}

			actual.Uniq()
			So(actual, ShouldResemble, expected)
		})

		Convey("Leaves values with the same name but unique city ids", func() {
			actual := CityNames{
				&CityName{Name: "Moscow", CityId: "1"},
				&CityName{Name: "Moscow", CityId: "2"},
			}

			expected := actual

			actual.Uniq()
			So(actual, ShouldResemble, expected)
		})
	})

	Convey("Search city names", t, func() {
		db := h.CreateDB(t)
		CreateCityNamesBucket(db)

		cityNames := CityNames{
			&CityName{
				Key: "moscow|2", Name: "Moscow", CityId: "2",
				Locale: "en", Population: 25000,
			},
			&CityName{
				Key: "moscow", Name: "Moscow", CityId: "1",
				Locale: "en", Population: 12000000,
			},
			&CityName{
				Key: "montreal", Name: "Montreal", CityId: "3",
				Locale: "en", Population: 1600000,
			},
			&CityName{
				Key: "moskau", Name: "Moskau", CityId: "1",
				Locale: "de", Population: 12000000,
			},
		}
		for _, cn := range cityNames {
			h.PutToBucket(t, db, CityNamesBucketName, cn.Key, cn.toString())
		}

		locales := []string{"ru", "en", "de"}
		result, err := searchCityNames(db, locales, "Mo", 5)

		Convey("Finds matching citynames without duplicates", func() {
			So(len(*result), ShouldEqual, 3)
			So((*result)[0], ShouldResemble, cityNames[1])
			So((*result)[1], ShouldResemble, cityNames[2])
			So((*result)[2], ShouldResemble, cityNames[0])
		})

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})
	})
}
