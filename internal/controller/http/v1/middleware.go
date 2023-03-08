package v1

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type AuthMiddleware struct {
	usecase usecase.Auth
}

func NewAuthMiddleware(uc usecase.Auth) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: uc,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err := m.usecase.ValidateToken(c.Request.Context(), headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		c.AbortWithStatus(status)
		return
	}

	// c.Set("user", user)

}

// Define custom middleware that skips the OpenTelemetry middleware for support routes
// func SkipOTLPMiddleware(c *gin.Context) {
// 	if c.Request.URL.Path == "/metrics" {
// 		c.Next()
// 		return
// 	}
// 	// } else if c.FullPath() == "/healthz" {
// 	// 	c.Next()
// 	// 	// return
// 	// } else if c.FullPath() == "/swagger/*any" {
// 	// 	c.Next()
// 	// 	// return
// 	// } else {
// 	// 	return
// 	// }
// 	otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME"))(c)
// }

// func SkipOTLPMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println(c.Request.URL.Path)
// 		if c.Request.URL.Path == "/metrics" {
// 			c.Next()
// 			return
// 		}
// 		if c.Request.URL.Path == "/healthz" {
// 			c.Next()
// 			return
// 		}
// 		if c.Request.URL.Path == "/secured/ping" {
// 			c.Next()
// 			return
// 		}
// 		if c.Request.URL.Path == "/swagger/*any" {
// 			c.Next()
// 			return
// 		}

// 		otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME"))(c)
// 	}
// }

// Define custom middleware that skips the OpenTelemetry middleware for support routes
func SkipOTLPMiddleware(c *gin.Context) {

	skipPaths := []string{
		"/metrics",
		"/healthz",
		"/secured/ping",
		"/swagger",
	}

	for _, path := range skipPaths {
		if strings.HasPrefix(c.Request.URL.Path, path) {
			// Skip the OpenTelemetry middleware for paths in skipPaths
			c.Next()
			return
		}
	}
	// Call the OpenTelemetry middleware for other routes
	otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME"))(c)
}
