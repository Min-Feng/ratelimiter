package main

import httpX "github.com/Min-Feng/ratelimiter/pkg/http"

func main() {
	router := httpX.NewRouter("8888")
	httpX.RegisterPath(router)
	router.QuickRun()
}
