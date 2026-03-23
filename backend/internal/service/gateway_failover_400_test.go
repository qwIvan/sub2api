package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldFailoverOn400(t *testing.T) {
	svc := &GatewayService{}

	tests := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "beta compatibility error",
			body: `{"error":{"message":"request requires beta feature anthropic-beta"}}`,
			want: true,
		},
		{
			name: "invalid model id keyword in message",
			body: `{"error":{"message":"INVALID_MODEL_ID"}}`,
			want: true,
		},
		{
			name: "improperly formed request in message",
			body: `{"error":{"message":"Improperly formed request"}}`,
			want: true,
		},
		{
			name: "non empty text blocks keyword in message",
			body: `{"error":{"message":"text content blocks must be non-empty"}}`,
			want: true,
		},
		{
			name: "generic bad request does not failover",
			body: `{"error":{"message":"invalid request payload"}}`,
			want: false,
		},
		{
			name: "empty body does not failover",
			body: `{}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, svc.shouldFailoverOn400([]byte(tt.body)))
		})
	}
}
