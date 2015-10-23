package main

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/cache"
	"github.com/lebedev-yury/cities/ds"
	"strings"
)

func warmUpSearchCache(db *bolt.DB, c *cache.Cache, locales []string, limit int) {
	enRunes := strings.Split(
		"a b c d e f g h i j k l m n o p q r s t u v w x y z", " ",
	)

	ruRunes := strings.Split(
		"а б в г д е ё ж з и й к л м н о п р с т у ф х ц ч ш щ ъ ы ь э ю я", " ",
	)

	jobs := make(chan string, 1000)
	results := make(chan string, 1000)

	go worker(jobs, results, db, c, locales, limit)
	jobsCount := 0

	for _, rune := range enRunes {
		jobsCount++
		jobs <- rune

		for _, secondRune := range enRunes {
			jobsCount++
			jobs <- rune + secondRune
		}
	}

	for _, rune := range ruRunes {
		jobsCount++
		jobs <- rune

		for _, secondRune := range ruRunes {
			jobsCount++
			jobs <- rune + secondRune
		}
	}

	close(jobs)

	for a := 1; a <= jobsCount; a++ {
		<-results
	}
}

func worker(
	jobs <-chan string, results chan<- string,
	db *bolt.DB, c *cache.Cache, locales []string, limit int,
) {
	for query := range jobs {
		ds.CachedCitiesSearch(db, c, locales, query, 5)
		results <- query
	}
}
