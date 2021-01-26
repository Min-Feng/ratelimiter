package main

import (
	"os"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/httpX"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func main() {
	cfg := configs.New("config")
	rateLimiter := limiter.New(&cfg, os.Getenv("LimiterKind"))
	router := httpX.NewRouter(&cfg, rateLimiter)
	httpX.RegisterPath(router)
	router.QuickRun()
}
