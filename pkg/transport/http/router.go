package http

import (
	"fluxio-backend/pkg/transport/http/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	r *gin.Engine
}

func NewRouter(authRoute *routes.AuthRoute) *Router {

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

	return &Router{
		r: router,
	}
}
