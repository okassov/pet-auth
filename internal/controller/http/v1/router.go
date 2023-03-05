package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
)

func NewRouter(handler *gin.Engine, a usecase.Auth) {

	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Healthcheck endpoint
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, a)
	}
}
