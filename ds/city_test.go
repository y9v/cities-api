package ds

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestCity(t *testing.T) {
	db := h.CreateDB(t)
	CreateCitiesBucket(db)
	CreateCountriesBucket(db)

	cityAttrs := []string{
		"Name",
		"DE",
		"10000000",
		"Latitude",
		"Longitude",
		"Timezone",
	}
	cityString := strings.Join(cityAttrs, "\t")

	Convey("City from string", t, func() {
		Convey("When the string is correct", func() {
			city, err := cityFromString("1", cityString)

			Convey("Sets the city id from param", func() {
				So(city.Id, ShouldEqual, "1")
			})

			Convey("Sets the city attributes from the string", func() {
				So(city.Name, ShouldEqual, cityAttrs[0])
				So(city.CountryCode, ShouldEqual, cityAttrs[1])
				So(city.Population, ShouldEqual, 10000000)
				So(city.Latitude, ShouldEqual, cityAttrs[3])
				So(city.Longitude, ShouldEqual, cityAttrs[4])
				So(city.Timezone, ShouldEqual, cityAttrs[5])
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the string is incorrect", func() {
			city, err := cityFromString("1", "")

			Convey("Leaves the city id blank", func() {
				So(city.Id, ShouldEqual, "")
			})

			Convey("Leaves the city attributes blank", func() {
				So(city.Name, ShouldEqual, "")
				So(city.CountryCode, ShouldEqual, "")
				So(city.Population, ShouldEqual, 0)
				So(city.Latitude, ShouldEqual, "")
				So(city.Longitude, ShouldEqual, "")
				So(city.Timezone, ShouldEqual, "")
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("City to string", t, func() {
		city := City{
			Id: "1", Name: "New York", CountryCode: "US", Population: 8600000,
			Latitude: "40.748817", Longitude: "-73.985428", Timezone: "USA/New York",
		}

		Convey("Joins the city properties with tab chars", func() {
			expected := "New York\tUS\t8600000\t40.748817\t-73.985428\tUSA/New York"
			So(expected, ShouldEqual, city.toString())
		})
	})

	Convey("Append city", t, func() {
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

	Convey("Find the city", t, func() {
		Convey("When the record exists", func() {
			h.PutToBucket(t, db, CitiesBucketName, "1", cityString)

			Convey("With includeCountry set to false", func() {
				city, err := FindCity(db, "1", false)

				Convey("Returns a city with attributes set", func() {
					expected, _ := cityFromString("1", cityString)
					So(city, ShouldResemble, expected)
				})

				Convey("Returns no error", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("With includeCountry set to true", func() {
				Convey("When country record exists", func() {
					countryAttrs := []string{"DE", "Germany", "en|Germany"}
					countryString := strings.Join(countryAttrs, "\t")
					h.PutToBucket(t, db, CountriesBucketName, "1", countryString)

					city, err := FindCity(db, "1", true)

					Convey("Returns a city with attributes set", func() {
						expected, _ := cityFromString("1", cityString)
						expected.Country, _ = countryFromString("1", countryString)
						So(city, ShouldResemble, expected)
					})

					Convey("Returns no error", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})

		Convey("When the record does not exist", func() {
			city, err := FindCity(db, "0", false)

			Convey("Returns a nil instead of a city", func() {
				So(city, ShouldBeNil)
			})

			Convey("Returns no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When an incorrect value is stored in the db", func() {
			h.PutToBucket(t, db, CitiesBucketName, "2", "")
			city, err := FindCity(db, "2", false)

			Convey("Returns an empty city", func() {
				So(city, ShouldResemble, &City{})
			})

			Convey("Returns an error", func() {
				So(err, ShouldNotBeNil)
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
