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

// Define custom middleware that skips the OpenTelemetry middleware for infra routes
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
