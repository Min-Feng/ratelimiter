package httpX

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Min-Feng/ratelimiter/pkg/configs"
	"github.com/Min-Feng/ratelimiter/pkg/limiter"
)

func NewDefaultRouter(cfg *configs.Config, limiter limiter.Limiter) *Router {
	r := NewRouter(cfg)
	limiterMiddleware := LimitIPAccessCountMiddleware(limiter)
	r.gin.Use(gin.Logger(), limiterMiddleware)
	return r
}

func NewRouter(cfg *configs.Config) *Router {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery())

	return &Router{
		gin:  router,
		addr: ":" + cfg.Port,
	}
}

type Router struct {
	gin  *gin.Engine
	addr string
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.gin.ServeHTTP(w, req)
}

func (r *Router) QuickRun() error {
	return r.gin.Run(r.addr)
}
