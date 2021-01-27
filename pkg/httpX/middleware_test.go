package httpX

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func Test_LimitIPAccessCount(t *testing.T) {
	cfg := configs.New("config")
	cfg.Port = "8168"
	cfg.Limiter.MaxLimitCount = 50
	cfg.Limiter.ResetCountIntervalSeconds = 3

	rateLimiter := limiter.New(&cfg, "local")
	router := NewRouter(&cfg, rateLimiter)
	RegisterPath(router)
	apiPath := "/hello"
	callCount := int(cfg.Limiter.MaxLimitCount + 10)

	for i := 1; i <= callCount; i++ {
		response, status := HTTPResponse(router, http.MethodGet, apiPath, nil)
		if int32(i) <= cfg.Limiter.MaxLimitCount {
			expectedCount := i
			ExpectedAccessEndpointOK200(t, response, status, expectedCount)
			continue
		}
		ExpectedAccessEndpointTooManyRequest429(t, response, status)
	}

	time.Sleep(cfg.Limiter.ResetCountInterval())
	response, status := HTTPResponse(router, http.MethodGet, apiPath, nil)
  expectedCount := 1
	ExpectedAccessEndpointOK200(t, response, status, expectedCount)
}

func ExpectedAccessEndpointOK200(t *testing.T, actualResponse string, actualStatus int, expectedCount int) {
	expectedResp := fmt.Sprintf("count: %v\nhello 192.0.2.1\n", expectedCount)
	assert.Equal(t, expectedResp, actualResponse)

	expectedStatus := http.StatusOK
	assert.Equal(t, expectedStatus, actualStatus)
}

func ExpectedAccessEndpointTooManyRequest429(t *testing.T, actualResponse string, actualStatus int) {
	expectedResp := "Error: ip=192.0.2.1 too many request\n"
	assert.Equal(t, expectedResp, actualResponse)

	expectedStatus := http.StatusTooManyRequests
	assert.Equal(t, expectedStatus, actualStatus)
}

func HTTPResponse(router http.Handler, httpMethod string, path string, body io.Reader) (respBody string, status int) {
	wRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(httpMethod, path, body)
	router.ServeHTTP(wRecorder, req)

	resp := wRecorder.Result()
  defer resp.Body.Close()
	return wRecorder.Body.String(), resp.StatusCode
}
