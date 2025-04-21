package http

import (
	"fluxio-backend/pkg/transport/http/routes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Address string
	Port    string
}

type Router struct {
	r *gin.Engine
}

func NewRouter(cfg RouterConfig, authRoute *routes.AuthRoute) *Router {

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1 := router.Group("/api/v1")
	{
		_ = v1.Group("/auth")
		{

		}
	}

	router.Run(fmt.Sprintf("%s:%s", cfg.Address, cfg.Port))

	return &Router{
		r: router,
	}
}
