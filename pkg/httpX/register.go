package httpX

import "time"

func RegisterPath(r *Router) {
	r.router.Use(LimitIPAccessCount(10, 10*time.Second))
	r.router.GET("/hello", Hello)
}
