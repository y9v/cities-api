package config

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestOptions(t *testing.T) {
	Convey("Load options from config file", t, func() {
		Convey("When config file does not exist", func() {
			defaultOptions := buildDefault()
			options := Load("foo.json")

			Convey("Returns default options", func() {
				So(options.Port, ShouldEqual, defaultOptions.Port)
				So(options.Timeout, ShouldEqual, defaultOptions.Timeout)
				So(options.CORSOrigins, ShouldResemble, defaultOptions.CORSOrigins)
				So(options.Locales, ShouldResemble, defaultOptions.Locales)
				So(options.MinPopulation, ShouldResemble, defaultOptions.MinPopulation)
				So(options.CitiesFile, ShouldEqual, defaultOptions.CitiesFile)
				So(options.CountriesFile, ShouldEqual, defaultOptions.CountriesFile)
				So(options.AlternateNamesFile, ShouldEqual, defaultOptions.AlternateNamesFile)
			})
		})

		Convey("When config file exists", func() {
			filename := h.CreateTempfile(
				t,
				`{
					 "Port": "3000",
					 "Timeout": 20,
					 "CORSOrigins": ["localhost", "example.com"],
					 "Locales": ["ru", "uk", "en"],
					 "MinPopulation": 10000,
					 "CountriesFile": "files/countries.txt",
					 "CitiesFile": "files/cities.txt",
					 "AlternateNamesFile": "files/alternate.txt"
				 }`,
			)

			options := Load(filename)

			Convey("Takes values from the config file", func() {
				So(options.Port, ShouldEqual, "3000")
				So(options.Timeout, ShouldEqual, 20)
				So(
					options.CORSOrigins, ShouldResemble,
					[]string{"localhost", "example.com"},
				)
				So(options.Locales, ShouldResemble, []string{"ru", "uk", "en"})
				So(options.MinPopulation, ShouldEqual, 10000)
				So(options.CountriesFile, ShouldEqual, "files/countries.txt")
				So(options.CitiesFile, ShouldEqual, "files/cities.txt")
				So(options.AlternateNamesFile, ShouldEqual, "files/alternate.txt")
			})
		})

		Convey("Panics when config file has wrong format", func() {
			filename := h.CreateTempfile(t, `"Port":}`)
			So(func() { Load(filename) }, ShouldPanic)
		})

		Convey("When the CONFIG env variable is set", func() {
			filename := h.CreateTempfile(t, `{"Port": "3001"}`)
			os.Setenv("CONFIG", filename)
			options := Load("foo.json")

			Convey("Options should be read from the env file", func() {
				So(options.Port, ShouldEqual, "3001")
			})
		})

		Convey("When the PORT env variable is set", func() {
			os.Setenv("PORT", "3030")
			options := Load("foo.json")

			Convey("Options should be read from the env file", func() {
				So(options.Port, ShouldEqual, "3030")
			})
		})
	})
}
