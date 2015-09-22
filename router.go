package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/lebedev-yury/cities/cache"
	"github.com/lebedev-yury/cities/config"
	"github.com/lebedev-yury/cities/middleware"
)

func newRouter(db *bolt.DB, options *config.Options, c *cache.Cache) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CommonHeaders(options.CORSOrigins))

	v1 := router.Group("/1.0")
	{
		v1.GET("/application/status", MakeApplicationStatusEndpoint(db))
		v1.GET("/cities/:id", MakeCityEndpoint(db))
		v1.GET("/search/cities", MakeSearchCitiesEndpoint(db, options, c))
	}

	return router
}
