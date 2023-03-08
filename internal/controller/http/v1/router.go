package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
)

func NewRouter(handler *gin.Engine, a usecase.Auth) {

	authMiddleware := NewAuthMiddleware(a)

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
	s := handler.Group("/secured", authMiddleware)
	{
		s.GET("ping", Ping)
	}

}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
