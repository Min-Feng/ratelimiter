package main

import (
	"fmt"
	"os"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/httpX"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func main() {
	cfg := configs.New("config")
	kind := os.Getenv("LimiterKind")
	if kind == "" {
		kind = "local"
	}
	fmt.Printf("limiter kind=%v\n", kind)

	rateLimiter := limiter.New(&cfg, kind)
	router := httpX.NewDefaultRouter(&cfg, rateLimiter)
	httpX.RegisterPath(router)

	fmt.Println("server start")
	router.QuickRun()
}
