package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// OpenAISessionHeaderCompatibility keeps the official Codex session-id header
// and the legacy underscore spelling aligned at the gateway boundary. The
// official hyphenated header wins when both spellings are present; downstream
// forwarding can therefore rely on session_id without losing the client value.
func OpenAISessionHeaderCompatibility() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c != nil && c.Request != nil {
			if sessionID := strings.TrimSpace(c.GetHeader("session-id")); sessionID != "" {
				c.Request.Header.Set("session_id", sessionID)
			}
		}
		c.Next()
	}
}
