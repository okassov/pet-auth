package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
)

type authRoutes struct {
	a usecase.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth) {
	r := &authRoutes{a}

	h := handler.Group("/auth")
	{
		h.GET("token")
	}
}
