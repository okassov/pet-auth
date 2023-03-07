package app

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/config"
	v1 "github.com/okassov/pet-auth/internal/controller/http/v1"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/internal/usecase/repository"
	"github.com/okassov/pet-auth/pkg/httpserver"
	"github.com/okassov/pet-auth/pkg/postgres"
)

func Run() {

	// Init Config
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// Repository
	pgConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.PG.PGUser,
		config.PG.PGPassword,
		config.PG.PGUrl,
		config.PG.PGPort,
		config.PG.PGDatabase)

	pg, err := postgres.New(pgConnString)
	if err != nil {
		fmt.Errorf("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()

	// Use case
	authUseCase := usecase.New(
		repository.New(pg),
		[]byte(config.JWT.JWTKey),
		time.Duration(config.JWT.JWTTtl),
	)

	handler := gin.New()

	// Router
	v1.NewRouter(handler, authUseCase)

	// Server
	serverPort := fmt.Sprintf(":%s", config.Server.Port)
	httpServer := httpserver.New(handler, serverPort)

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		fmt.Println("Error server start")
	}
}
