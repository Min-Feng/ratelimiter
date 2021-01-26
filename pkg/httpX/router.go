package httpX

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func NewRouter(cfg *configs.Config, limiter limiter.Limiter) *Router {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()

	limiterMiddleware := LimitIPAccessCountMiddleware(limiter)
	router.Use(gin.Recovery(), gin.Logger(), limiterMiddleware)

	return &Router{
		router: router,
		addr:   ":" + cfg.Port,
	}
}

type Router struct {
	router *gin.Engine
	addr   string
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) QuickRun() error {
	return r.router.Run(r.addr)
}
