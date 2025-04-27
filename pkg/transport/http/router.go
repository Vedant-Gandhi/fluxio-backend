package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Address string
	Port    string
}

// RouteRegistrar interface that all route groups should implement
type RouteRegistrar interface {
	RegisterRoutes(router *gin.Engine)
}

type Router struct {
	engine  *gin.Engine
	address string
	port    string
}

func NewRouter(cfg RouterConfig, routeRegistrars ...RouteRegistrar) *Router {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Register all routes provided by the route registrars
	for _, registrar := range routeRegistrars {
		registrar.RegisterRoutes(router)
	}

	return &Router{
		engine:  router,
		address: cfg.Address,
		port:    cfg.Port,
	}
}

// Start starts the HTTP server
func (r *Router) Start() error {
	return r.engine.Run(fmt.Sprintf("%s:%s", r.address, r.port))
}
