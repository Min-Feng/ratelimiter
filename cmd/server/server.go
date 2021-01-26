package main

import (
	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/httpX"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func main() {
	cfg := configs.New("config")
	rateLimiter := limiter.NewLocalLimiter(cfg.Limiter.MaxLimitCount, cfg.Limiter.ResetCountInterval())
	router := httpX.NewRouter(&cfg, rateLimiter)
	httpX.RegisterPath(router)
	router.QuickRun()
}
