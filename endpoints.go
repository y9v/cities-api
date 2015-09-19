package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/ds"
)

func MakeApplicationStatusEndpoint(db *bolt.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		appStatus := ds.GetAppStatus(db)
		c.JSON(200, appStatus)
	}
}

func MakeCityEndpoint(db *bolt.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		city, err := ds.FindCity(db, c.Param("id"), true)

		if err != nil {
			c.JSON(500, nil)
		} else if city == nil {
			c.JSON(404, nil)
		} else {
			c.JSON(200, city)
		}
	}
}

func MakeSearchCitiesEndpoint(
	db *bolt.DB, options *config.Options,
) func(*gin.Context) {
	return func(c *gin.Context) {
		query := c.Query("q")

		if query == "" {
			c.JSON(200, nil)
		} else {
			cities, err := ds.SearchCities(db, options.Locales, query, 5)

			if err != nil {
				c.JSON(500, nil)
			} else {
				c.JSON(200, cities)
			}
		}
	}
}
