package main

import (
	"github.com/Min-Feng/ratelimiter/pkg/configs"
	httpX "github.com/Min-Feng/ratelimiter/pkg/httpX"
)

func main() {
	cfg := configs.New("config")
	// spew.Dump(cfg)

	router := httpX.NewRouter(cfg.Port)
	httpX.RegisterPath(cfg.Limiter, router)
	router.QuickRun()
}
