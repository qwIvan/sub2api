package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOpenAISessionHeaderCompatibilityPrefersOfficialHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(OpenAISessionHeaderCompatibility())
	router.POST("/", func(c *gin.Context) {
		if got := c.GetHeader("session_id"); got != "official" {
			t.Fatalf("session_id = %q, want official", got)
		}
		if got := c.GetHeader("session-id"); got != "official" {
			t.Fatalf("session-id = %q, want official", got)
		}
	})

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("session-id", "official")
	req.Header.Set("session_id", "legacy")
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestOpenAISessionHeaderCompatibilityLeavesLegacyHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(OpenAISessionHeaderCompatibility())
	router.POST("/", func(c *gin.Context) {
		if got := c.GetHeader("session_id"); got != "legacy" {
			t.Fatalf("session_id = %q, want legacy", got)
		}
	})

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("session_id", "legacy")
	router.ServeHTTP(httptest.NewRecorder(), req)
}
