package middleware

import (
	"github.com/drewolson/testflight"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestCommonHeaders(t *testing.T) {
	Convey("Test common headers", t, func() {
		corsOrigins := []string{"http://localhost", "http://127.0.0.1"}

		testflight.WithServer(
			testHandler(corsOrigins),
			func(r *testflight.Requester) {
				headers := r.Get("/ping").RawResponse.Header

				Convey("Sets content-type header", func() {
					actual := headers["Content-Type"]
					So(actual, ShouldResemble, []string{"application/json"})
				})

				Convey("Sets CORS header", func() {
					actual := headers["Access-Control-Allow-Origin"]
					So(actual, ShouldResemble, []string{strings.Join(corsOrigins, ",")})
				})

				Convey("Sets access-control-allow header", func() {
					actual := headers["Access-Control-Allow-Headers"]
					So(actual, ShouldResemble, []string{
						"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
					})
				})
			},
		)
	})
}

func testHandler(corsOrigins []string) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(CommonHeaders(corsOrigins))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r
}
