package httpX

func RegisterPath(r *Router) {
	r.router.GET("/hello", Hello)
}
