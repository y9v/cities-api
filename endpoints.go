package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/lebedev-yury/cities/cache"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/ds"
)

func MakeApplicationStatusEndpoint(db *bolt.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		appStatus := ds.GetAppStatus(db)
		c.JSON(200, appStatus)
	}
}

func MakeCityEndpoint(db *bolt.DB, options *config.Options) func(*gin.Context) {
	return func(c *gin.Context) {
		city, err := ds.FindCity(db, c.Param("id"), true)

		if err != nil {
			c.JSON(500, nil)
		} else if city == nil {
			c.JSON(404, nil)
		} else {
			c.JSON(200, city.ForSerialization(c.Request.URL, options))
		}
	}
}

func MakeCitiesEndpoint(
	db *bolt.DB, options *config.Options, cache *cache.Cache,
) func(*gin.Context) {
	return func(c *gin.Context) {
		query := ds.PrepareCityNameKey(c.Query("q"))

		if query == "" {
			c.JSON(200, nil)
		} else {
			cities, err := ds.CachedCitiesSearch(db, cache, options.Locales, query, 5)

			if err != nil {
				c.JSON(500, nil)
			} else {
				c.JSON(200, cities.ForSerialization(c.Request.URL, options))
			}
		}
	}
}
