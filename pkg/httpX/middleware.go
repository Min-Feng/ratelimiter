package httpX

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func LimitIPAccessCountMiddleware(rateLimiter limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		count, err := rateLimiter.Allow(ip)
		if err != nil {
			defer c.Abort()

			if errors.Is(err, limiter.ErrExceedMaxCount) {
				c.Data(
					http.StatusTooManyRequests,
					"text/plain",
					[]byte(fmt.Sprintf("Error: ip=%v too many request\n", ip)),
				)
				return
			}

			c.Data(http.StatusInternalServerError, "text/plain", []byte(fmt.Sprintf("Error: %v\n", err)))
			return
		}

		c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("count: %v\n", count)))
	}
}
