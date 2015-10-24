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
	corsOrigins := []string{"http://example.com"}

	Convey("Test common headers", t, func() {
		testflight.WithServer(
			testHandler(corsOrigins),
			func(r *testflight.Requester) {
				headers := r.Get("/ping").RawResponse.Header

				Convey("Sets content-type header", func() {
					actual := headers.Get("Content-Type")
					So(actual, ShouldEqual, "application/vnd.api+json")
				})

				Convey("Sets access-control-allow header", func() {
					actual := headers.Get("Access-Control-Allow-Headers")
					So(actual, ShouldEqual, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				})
			},
		)
	})

	Convey("Test CORS headers", t, func() {
		testflight.WithServer(
			testHandler(corsOrigins),
			func(r *testflight.Requester) {
				headerName := "Access-Control-Allow-Origin"
				request, _ := http.NewRequest(
					"GET", "/ping", strings.NewReader(""),
				)

				Convey("Returns CORS header for allowed origin", func() {
					request.Header.Add("Origin", "http://example.com")
					actual := r.Do(request).RawResponse.Header.Get(headerName)
					So(actual, ShouldEqual, "http://example.com")
				})

				Convey("Does not return CORS header for not allowed origin", func() {
					request.Header.Add("Origin", "http://not-allowed.com")
					actual := r.Do(request).RawResponse.Header.Get(headerName)
					So(actual, ShouldEqual, "")
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
