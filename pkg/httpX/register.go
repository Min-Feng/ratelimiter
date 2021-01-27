package httpX

func RegisterPath(r *Router) {
	r.gin.GET("/hello", Hello)
}
