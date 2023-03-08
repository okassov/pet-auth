package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/config"
	v1 "github.com/okassov/pet-auth/internal/controller/http/v1"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/internal/usecase/repository"
	"github.com/okassov/pet-auth/pkg/httpserver"
	"github.com/okassov/pet-auth/pkg/logger"
	"github.com/okassov/pet-auth/pkg/postgres"
	"github.com/okassov/pet-auth/pkg/tracer"
)

func Run(config config.Config) {

	// Init Logger
	l := logger.New()

	// Init TracerProvider
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Check OTEL_SERVICE_NAME env (required)
	_, ok := os.LookupEnv("OTEL_SERVICE_NAME")
	if !ok {
		log.Fatal("Cannot find OTEL_SERVICE_NAME environment variable.")
	}
	err = tracer.NewTracerProvider(f)
	if err != nil {
		log.Fatal(err)
	}

	// Repository
	pgConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.PG.PGUser,
		config.PG.PGPassword,
		config.PG.PGUrl,
		config.PG.PGPort,
		config.PG.PGDatabase)

	pg, err := postgres.New(
		pgConnString,
		postgres.MaxPoolSize(config.PG.PGMaxPool),
		postgres.ConnAttempts(config.PG.PGConnAttempts),
		postgres.ConnTimeout(time.Duration(config.PG.PGConnTimeout)),
	)
	if err != nil {
		l.Error("app - Run - postgres.New: %w", err)
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
	v1.NewRouter(handler, authUseCase, l)

	// Server
	httpServer := httpserver.New(handler, httpserver.Port(config.Server.Port))

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown", err)
	}
}
