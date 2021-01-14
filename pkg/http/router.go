package httpX

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(port string) *Router {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery())

	return &Router{
		router: router,
		addr:   ":" + port,
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
