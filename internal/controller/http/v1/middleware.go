package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
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

	user, err := m.usecase.ValidateToken(c.Request.Context(), headerParts[1])
	if err != nil {
		status := http.StatusInternalServerError
		c.AbortWithStatus(status)
		return
	}

	c.Set("user", user)

}
