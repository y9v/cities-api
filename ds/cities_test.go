package ds

import (
	"github.com/lebedev-yury/cities/cache"
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestCities(t *testing.T) {
	Convey("Append city", t, func() {
		db := h.CreateDB(t)
		CreateCountriesBucket(db)

		h.PutToBucket(
			t, db, CountriesBucketName, "1",
			"US\tUnited States\ten|United States;ru|Соединенные Штаты",
		)

		cities := Cities{
			Cities: []*City{
				&City{Name: "Venice"}, &City{Name: "Moscow"},
			},
		}

		Convey("When no city with the same name is in the collection", func() {
			city := City{Name: "London"}
			actual := appendCity(db, cities.Cities, &city, "en")

			Convey("Adds the city to the array", func() {
				So(len(actual), ShouldEqual, 3)
			})

			Convey("Leaves the city name unchanged", func() {
				So(actual[2].Name, ShouldEqual, city.Name)
			})
		})

		Convey("When city with the same name is in the collection", func() {
			Convey("When the country exists for the given code", func() {
				city := City{Name: "Venice", CountryCode: "US"}

				Convey("Default locale", func() {
					actual := appendCity(db, cities.Cities, &city, "en")

					Convey("Adds the city to the array", func() {
						So(len(actual), ShouldEqual, 3)
					})

					Convey("Adds the country name to the city name", func() {
						So(actual[2].Name, ShouldEqual, "Venice, United States")
					})
				})

				Convey("Some other locale", func() {
					actual := appendCity(db, cities.Cities, &city, "ru")

					Convey("Adds the city to the array", func() {
						So(len(actual), ShouldEqual, 3)
					})

					Convey("Adds the country name to the city name", func() {
						So(actual[2].Name, ShouldEqual, "Venice, Соединенные Штаты")
					})
				})
			})

			Convey("When no country exists for the given code", func() {
				city := City{Name: "Moscow", CountryCode: "MO"}
				actual := appendCity(db, cities.Cities, &city, "en")

				Convey("Doesn't adds the city to the array", func() {
					So(len(actual), ShouldEqual, 2)
				})
			})
		})
	})

	Convey("Search cites", t, func() {
		db := h.CreateDB(t)
		CreateCitiesBucket(db)
		CreateCityNamesBucket(db)

		cityNames := CityNames{
			&CityName{
				Key: "montreal", Name: "Montréal", CityId: 1,
				Locale: "fr", Population: 1600000,
			},
			&CityName{
				Key: "moscow", Name: "Moskau", CityId: 2,
				Locale: "de", Population: 12000000,
			},
		}
		for _, cn := range cityNames {
			h.PutToBucket(t, db, CityNamesBucketName, cn.Key, cn.toString())
		}

		cities := []*City{
			&City{ID: 1, Name: "Montreal"},
			&City{ID: 2, Name: "Moscow"},
		}
		for _, city := range cities {
			h.PutToBucket(
				t, db, CitiesBucketName, strconv.Itoa(city.ID), city.toString(),
			)
		}

		locales := []string{"ru", "en", "de"}

		Convey("Non cached search", func() {
			result, err := SearchCities(db, locales, "Mo", 5)

			Convey("Finds matching cities", func() {
				So(len(result.Cities), ShouldEqual, 2)
				So(result.Cities[0].ID, ShouldEqual, cities[1].ID)
				So(result.Cities[1].ID, ShouldEqual, cities[0].ID)
			})

			Convey("Sets the city names from the mathing cityname", func() {
				So(result.Cities[0].Name, ShouldEqual, cityNames[1].Name)
				So(result.Cities[1].Name, ShouldEqual, cityNames[0].Name)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Cached search", func() {
			c := cache.New()

			Convey("For short queries", func() {
				results, expectedErr := SearchCities(db, locales, "Mo", 5)
				cacheMissResults, actualErr := CachedCitiesSearch(db, c, locales, "Mo", 5)
				cacheHitResults, _ := CachedCitiesSearch(db, c, locales, "Mo", 5)

				Convey("Returns results from search if cache miss", func() {
					So(cacheMissResults, ShouldResemble, results)
					So(actualErr, ShouldEqual, expectedErr)
				})

				Convey("Returns results from cache if cache hit", func() {
					So(cacheHitResults, ShouldResemble, cacheMissResults)
				})
			})

			Convey("For longer queries", func() {
				expected, expectedErr := SearchCities(db, locales, "Moscow", 5)
				actual, actualErr := CachedCitiesSearch(db, c, locales, "Moscow", 5)

				Convey("Returns results from search if cache miss", func() {
					So(actual, ShouldResemble, expected)
					So(actualErr, ShouldEqual, expectedErr)
				})
			})
		})
	})
}
