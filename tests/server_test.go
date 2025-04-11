package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// Setup your router, just like you did in your main function, and register your routes
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Perform a GET request with that router (using `http` package's `NewRequest` and `httptest`'s `NewRecorder`)
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, resp.Code, http.StatusOK)
	// Assert the response body is what we expect.
	assert.Equal(t, resp.Body.String(), "pong")
}
