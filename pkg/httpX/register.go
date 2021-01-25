package httpX

import (
	"github.com/Min-Feng/ratelimiter/pkg/configs"
)

func RegisterPath(cfg configs.Limiter, r *Router) {
	rateLimiterMiddleware := LimitIPAccessCount(cfg.MaxLimitCount, cfg.ResetCountInterval())
	r.router.Use(rateLimiterMiddleware)
	r.router.GET("/hello", Hello)
}
