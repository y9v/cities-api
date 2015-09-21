package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
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
				Key: "montreal", Name: "Montréal", CityId: "1",
				Locale: "fr", Population: 1600000,
			},
			&CityName{
				Key: "moscow", Name: "Moskau", CityId: "2",
				Locale: "de", Population: 12000000,
			},
		}
		for _, cn := range cityNames {
			h.PutToBucket(t, db, CityNamesBucketName, cn.Key, cn.toString())
		}

		cities := []*City{
			&City{Id: "1", Name: "Montreal"},
			&City{Id: "2", Name: "Moscow"},
		}
		for _, city := range cities {
			h.PutToBucket(t, db, CitiesBucketName, city.Id, city.toString())
		}

		locales := []string{"ru", "en", "de"}
		result, err := SearchCities(db, locales, "Mo", 5)

		Convey("Finds matching cities", func() {
			So(len(result.Cities), ShouldEqual, 2)
			So(result.Cities[0].Id, ShouldEqual, cities[1].Id)
			So(result.Cities[1].Id, ShouldEqual, cities[0].Id)
		})

		Convey("Sets the city names from the mathing cityname", func() {
			So(result.Cities[0].Name, ShouldEqual, cityNames[1].Name)
			So(result.Cities[1].Name, ShouldEqual, cityNames[0].Name)
		})

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})
	})
}
