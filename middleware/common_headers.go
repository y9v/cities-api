package middleware

import (
	"github.com/gin-gonic/gin"
)

func CommonHeaders(corsOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, domain := range corsOrigins {
			if c.Request.Header.Get("Origin") == domain {
				c.Writer.Header().Set("Access-Control-Allow-Origin", domain)
			}
		}

		c.Writer.Header().Set("Content-Type", "application/vnd.api+json")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)

		c.Next()
	}
}
