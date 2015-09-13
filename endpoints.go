package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func applicationStatusEndpoint(c *gin.Context) {
	appStatus := GetAppStatus()
	c.JSON(200, appStatus)
}

func cityEndpoint(c *gin.Context) {
	city, err := FindCity(c.Param("id"))

	if err != nil {
		c.JSON(500, nil)
	} else if city == nil {
		c.JSON(404, nil)
	} else {
		c.JSON(200, city)
	}
}

func searchCitiesEndpoint(c *gin.Context) {
	query := strings.Trim(strings.ToLower(c.Query("q")), "| ")

	if query == "" {
		c.JSON(200, nil)
	} else {
		cities, err := SearchCitiesByCityName(query, 5)

		if err != nil {
			c.JSON(500, nil)
		} else {
			c.JSON(200, cities)
		}
	}
}
