package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAICacheMaskProbabilityNumerator(t *testing.T) {
	require.Equal(t, 50, openAICacheMaskProbabilityNumerator("admin"))
	require.Equal(t, 10, openAICacheMaskProbabilityNumerator("yangyue"))
	require.Equal(t, 0, openAICacheMaskProbabilityNumerator(""))
	require.Equal(t, 0, openAICacheMaskProbabilityNumerator("user"))
}

func TestShouldMaskOpenAICacheUsageForClientOnlyConfiguredUsersParticipate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	newContext := func(userid string) *gin.Context {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)
		context.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		if userid != "" {
			context.Request.Header.Set("Auth-Userid", userid)
		}
		return context
	}

	openAIAccount := &Account{Platform: PlatformOpenAI}
	require.False(t, shouldMaskOpenAICacheUsageForClient(openAIAccount, newContext("")))
	require.False(t, shouldMaskOpenAICacheUsageForClient(openAIAccount, newContext("user")))
	require.False(t, shouldMaskOpenAICacheUsageForClient(&Account{Platform: PlatformAnthropic}, newContext("admin")))
}

func TestMaskOpenAICacheUsageForClient(t *testing.T) {
	body := []byte(`{"usage":{"input_tokens_details":{"cached_tokens":7},"prompt_tokens_details":{"cached_tokens":8},"cache_read_input_tokens":9,"cached_tokens":10,"prompt_cache_hit_tokens":11},"response":{"usage":{"input_tokens_details":{"cached_tokens":12}}},"data":{"response":{"usage":{"cache_read_input_tokens":13}}},"message":"keep"}`)

	masked, changed := maskOpenAICacheUsageForClient(body)

	require.True(t, changed)
	require.Equal(t, int64(0), gjson.GetBytes(masked, "usage.input_tokens_details.cached_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "usage.prompt_tokens_details.cached_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "usage.cache_read_input_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "usage.cached_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "usage.prompt_cache_hit_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "response.usage.input_tokens_details.cached_tokens").Int())
	require.Equal(t, int64(0), gjson.GetBytes(masked, "data.response.usage.cache_read_input_tokens").Int())
	require.Equal(t, "keep", gjson.GetBytes(masked, "message").String())
}

func TestMaskOpenAICacheUsageForClientWithoutCacheFields(t *testing.T) {
	body := []byte(`{"usage":{"input_tokens":100,"output_tokens":20},"message":"keep"}`)

	masked, changed := maskOpenAICacheUsageForClient(body)

	require.False(t, changed)
	require.Equal(t, string(body), string(masked))
}
