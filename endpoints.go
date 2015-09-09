package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

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
	query := strings.Replace(strings.ToLower(c.Query("q")), "|", "", -1)
	cities, err := SearchCitiesByCityName(query, 5)

	if err != nil {
		c.JSON(500, nil)
	} else {
		c.JSON(200, cities)
	}
}
