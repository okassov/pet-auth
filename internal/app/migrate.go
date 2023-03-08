//go:build migrate

package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}

	databasePort, ok := os.LookupEnv("PG_PORT")
	if !ok || len(databasePort) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_PORT")
	}

	databaseUser, ok := os.LookupEnv("PG_USER")
	if !ok || len(databaseUser) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_USER")
	}

	databasePassword, ok := os.LookupEnv("PG_PASSWORD")
	if !ok || len(databasePassword) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_PASSWORD")
	}

	databaseName, ok := os.LookupEnv("PG_DATABASE")
	if !ok || len(databaseName) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_DATABASE")
	}

	databaseConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		databaseUser,
		databasePassword,
		databaseURL,
		databasePort,
		databaseName)

	databaseConnString += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseConnString)

		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
