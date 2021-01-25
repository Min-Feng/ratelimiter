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
)

func Test_LimitIPAccessCount(t *testing.T) {
	cfg := configs.Config{
		Port: "8888",
		Limiter: configs.Limiter{
			MaxLimitCount:             10,
			ResetCountIntervalSeconds: 2,
			RemoveIntervalHour:        1,
		},
	}

	router := NewRouter(cfg.Port)
	RegisterPath(cfg.Limiter, router)
	path := "/hello"
	callCount := 20

	for i := 1; i <= callCount; i++ {
		response, status := HTTPResponse(router, http.MethodGet, path, nil)

		if int32(i) <= cfg.Limiter.MaxLimitCount {
			expectedResp := fmt.Sprintf("count: %v\nhello 192.0.2.1\n", i)
			assert.Equal(t, expectedResp, response)

			expectedStatus := http.StatusOK
			assert.Equal(t, expectedStatus, status)
			continue
		}

		expectedResp := "Error: ip=192.0.2.1 too many request\n"
		assert.Equal(t, expectedResp, response)

		expectedStatus := http.StatusTooManyRequests
		assert.Equal(t, expectedStatus, status)
	}

	time.Sleep(cfg.Limiter.ResetCountInterval())
	response, status := HTTPResponse(router, http.MethodGet, path, nil)

	expectedResp := fmt.Sprintf("count: %v\nhello 192.0.2.1\n", 1)
	assert.Equal(t, expectedResp, response)

	expectedStatus := http.StatusOK
	assert.Equal(t, expectedStatus, status)
}

func HTTPResponse(router http.Handler, httpMethod string, url string, body io.Reader) (resp string, status int) {
	wRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(httpMethod, url, body)
	router.ServeHTTP(wRecorder, req)
	actualBody := string(wRecorder.Body.Bytes())
	return actualBody, wRecorder.Result().StatusCode
}
