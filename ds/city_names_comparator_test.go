package ds

import (
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestCityNamesComparator(t *testing.T) {
	Convey("City names comparator", t, func() {
		Convey("Sorts by population", func() {
			actual := CityNames{
				&CityName{Name: "Foo", Population: 10000},
				&CityName{Name: "Bar", Population: 12000000},
				&CityName{Name: "Baz", Population: 500000},
			}

			expected := CityNames{
				&CityName{Name: "Bar", Population: 12000000},
				&CityName{Name: "Baz", Population: 500000},
				&CityName{Name: "Foo", Population: 10000},
			}

			sort.Sort(CityNamesComparator{actual, []string{"en"}})
			So(actual, ShouldResemble, expected)
		})

		Convey("Sorts by locale when the population is the same", func() {
			actual := CityNames{
				&CityName{Name: "Foo", Locale: "de", Population: 10000},
				&CityName{Name: "Bar", Locale: "ru", Population: 10000},
				&CityName{Name: "Baz", Locale: "", Population: 10000},
				&CityName{Name: "Qux", Locale: "", Population: 10000},
				&CityName{Name: "Norf", Locale: "en", Population: 10000},
			}

			expected := CityNames{
				&CityName{Name: "Bar", Locale: "ru", Population: 10000},
				&CityName{Name: "Norf", Locale: "en", Population: 10000},
				&CityName{Name: "Foo", Locale: "de", Population: 10000},
				&CityName{Name: "Baz", Locale: "", Population: 10000},
				&CityName{Name: "Qux", Locale: "", Population: 10000},
			}

			locales := []string{"ru", "en", "de"}

			sort.Sort(CityNamesComparator{actual, locales})
			So(actual, ShouldResemble, expected)
		})
	})
}
