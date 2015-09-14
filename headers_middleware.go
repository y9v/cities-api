package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func HeadersMiddleware() gin.HandlerFunc {
	corsHeader := strings.Join(options.CORSOrigins, ",")

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", corsHeader)
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)

		c.Next()
	}
}
