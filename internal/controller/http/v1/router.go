package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
)

func NewRouter(handler *gin.Engine, a usecase.Auth) {

	authMiddleware := NewAuthMiddleware(a)
	// api := router.Group("/api", authMiddleware)

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
	t := handler.Group("/test", authMiddleware)
	{
		t.GET("ping", testPing)
	}

}

func testPing(c *gin.Context) {
	c.Status(http.StatusOK)
}
