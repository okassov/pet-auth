package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/pkg/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/okassov/pet-auth/docs"
)

//	@title			Pet Auth Service
//	@description	Using a authentication service
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/v1
func NewRouter(handler *gin.Engine, a usecase.Auth, l logger.LoggerInterface) {

	// authMiddleware := NewAuthMiddleware(a)

	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Healthcheck endpoint
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, a, l)
	}
	// s := handler.Group("/secured", authMiddleware)
	// {
	// 	s.GET("ping", Ping)
	// }

}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
