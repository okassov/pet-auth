package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/pkg/httpserver"
)

func Run() {

	// Use case
	authUseCase := usecase.New()

	handler := gin.New()

	v1.NewRouter(handler, authUseCase)

	httpServer := httpserver.New(handler)

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		fmt.Println("Error server start")
	}
}
