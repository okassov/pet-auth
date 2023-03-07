package app

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/okassov/pet-auth/internal/controller/http/v1"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/internal/usecase/repository"
	"github.com/okassov/pet-auth/pkg/httpserver"
	"github.com/okassov/pet-auth/pkg/postgres"
)

var pg_url = "postgres://auth:auth@localhost:5432/auth"
var signing_key = "signing_key"
var token_ttl = 86400

func Run() {

	// Repository
	pg, err := postgres.New(pg_url)
	if err != nil {
		fmt.Errorf("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()

	// Use case
	authUseCase := usecase.New(
		repository.New(pg),
		[]byte(signing_key),
		time.Duration(token_ttl),
	)

	handler := gin.New()

	v1.NewRouter(handler, authUseCase)

	httpServer := httpserver.New(handler)

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		fmt.Println("Error server start")
	}
}
