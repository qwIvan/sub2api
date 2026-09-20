package repository

import (
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// mirrorOpenAISessionHeader applies the final outbound compatibility projection
// after all gateway/account isolation logic has run. The underscore spelling is
// the value produced by the existing OpenAI builders; expose that same value
// under Codex's official hyphenated header without changing other upstreams.
func mirrorOpenAISessionHeader(req *http.Request) {
	if req == nil || req.Header == nil || service.HTTPUpstreamProfileFromContext(req.Context()) != service.HTTPUpstreamProfileOpenAI {
		return
	}
	if sessionID := strings.TrimSpace(req.Header.Get("session_id")); sessionID != "" {
		req.Header.Set("session-id", sessionID)
	}
}
