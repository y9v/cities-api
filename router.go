package main

import (
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	router := gin.Default()
	router.Use(HeadersMiddleware())

	v1 := router.Group("/1.0")
	{
		v1.GET("/application/status", applicationStatusEndpoint)
		v1.GET("/cities/:id", cityEndpoint)
		v1.GET("/search/cities", searchCitiesEndpoint)
	}

	return router
}
