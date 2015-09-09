package main

import (
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/1.0")
	{
		v1.GET("/cities.json/:id", cityEndpoint)
		v1.GET("/search/cities.json", searchCitiesEndpoint)
	}

	return router
}
