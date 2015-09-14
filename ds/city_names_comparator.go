package ds

type CityNamesComparator struct {
	CityNames
	Locales []string
}

func (slice CityNames) Len() int {
	return len(slice)
}

func (slice CityNames) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice CityNamesComparator) Less(i, j int) bool {
	if slice.CityNames[i].Population == slice.CityNames[j].Population {
		for _, locale := range slice.Locales {
			if slice.CityNames[i].Locale == locale {
				return true
			} else if slice.CityNames[j].Locale == locale {
				return false
			}
		}
		return false
	} else {
		return slice.CityNames[i].Population > slice.CityNames[j].Population
	}
}
