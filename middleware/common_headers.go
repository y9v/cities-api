package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func CommonHeaders(corsOrigins []string) gin.HandlerFunc {
	corsHeader := strings.Join(corsOrigins, ",")

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
