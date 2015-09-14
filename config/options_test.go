package config

import (
	h "github.com/lebedev-yury/cities/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestOptions(t *testing.T) {
	Convey("Load options from config file", t, func() {
		options := Options{
			Port:               "80",
			Timeout:            10,
			CORSOrigins:        []string{"127.0.0.1"},
			Locales:            []string{"en", "de"},
			CitiesFile:         "~/cities.txt",
			AlternateNamesFile: "~/alternate.txt",
		}

		Convey("When config file does not exist", func() {
			Load(&options, "foo.json")

			Convey("Keeps original values", func() {
				So(options.Port, ShouldEqual, "80")
				So(options.Timeout, ShouldEqual, 10)
				So(options.CORSOrigins, ShouldResemble, []string{"127.0.0.1"})
				So(options.Locales, ShouldResemble, []string{"en", "de"})
				So(options.CitiesFile, ShouldEqual, "~/cities.txt")
				So(options.AlternateNamesFile, ShouldEqual, "~/alternate.txt")
			})
		})

		Convey("When config file exists", func() {
			filename := h.CreateTempfile(
				`{
					 "Port": "3000",
					 "Timeout": 20,
					 "CORSOrigins": ["localhost", "example.com"],
					 "Locales": ["ru", "uk", "en"],
					 "CitiesFile": "files/cities.txt",
					 "AlternateNamesFile": "files/alternate.txt"
				 }`, t,
			)

			Load(&options, filename)

			Convey("Takes values from the config file", func() {
				So(options.Port, ShouldEqual, "3000")
				So(options.Timeout, ShouldEqual, 20)
				So(
					options.CORSOrigins, ShouldResemble,
					[]string{"localhost", "example.com"},
				)
				So(options.Locales, ShouldResemble, []string{"ru", "uk", "en"})
				So(options.CitiesFile, ShouldEqual, "files/cities.txt")
				So(options.AlternateNamesFile, ShouldEqual, "files/alternate.txt")
			})
		})

		Convey("Panics when config file has wrong format", func() {
			filename := h.CreateTempfile(`"Port":}`, t)
			So(func() { Load(&options, filename) }, ShouldPanic)
		})

		Convey("When the CONFIG env variable is set", func() {
			filename := h.CreateTempfile(`{"Port": "3001"}`, t)
			os.Setenv("CONFIG", filename)
			Load(&options, "foo.json")

			Convey("Options should be read from the env file", func() {
				So(options.Port, ShouldEqual, "3001")
			})
		})

		Convey("When the PORT env variable is set", func() {
			os.Setenv("PORT", "3030")
			Load(&options, "foo.json")

			Convey("Options should be read from the env file", func() {
				So(options.Port, ShouldEqual, "3030")
			})
		})
	})
}
