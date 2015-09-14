package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lebedev-yury/cities/ds"
)

func applicationStatusEndpoint(c *gin.Context) {
	appStatus := GetAppStatus()
	c.JSON(200, appStatus)
}

func cityEndpoint(c *gin.Context) {
	city, err := ds.FindCity(db, c.Param("id"))

	if err != nil {
		c.JSON(500, nil)
	} else if city == nil {
		c.JSON(404, nil)
	} else {
		c.JSON(200, city)
	}
}

func searchCitiesEndpoint(c *gin.Context) {
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
