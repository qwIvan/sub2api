package repository

import (
	"context"
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestMirrorOpenAISessionHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "https://example.test/v1/responses", nil)
	req = req.WithContext(service.WithHTTPUpstreamProfile(context.Background(), service.HTTPUpstreamProfileOpenAI))
	req.Header.Set("session_id", "isolated-session")
	req.Header.Set("session-id", "stale-session")

	mirrorOpenAISessionHeader(req)
	if got := req.Header.Get("session-id"); got != "isolated-session" {
		t.Fatalf("session-id = %q, want isolated-session", got)
	}
}

func TestMirrorOpenAISessionHeaderDoesNotTouchOtherProfiles(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "https://example.test/v1/messages", nil)
	req.Header.Set("session_id", "legacy-session")

	mirrorOpenAISessionHeader(req)
	if got := req.Header.Get("session-id"); got != "" {
		t.Fatalf("session-id = %q, want empty for non-OpenAI profile", got)
	}
}
