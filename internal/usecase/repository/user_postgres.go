package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/okassov/pet-auth/internal/entity"
	"github.com/okassov/pet-auth/pkg/postgres"
	"go.opentelemetry.io/otel"
)

type UserRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, a entity.User) error {

	// Tracer
	tracerName := os.Getenv("OTEL_SERVICE_NAME")
	_, span := otel.GetTracerProvider().Tracer(tracerName).Start(ctx, "RepositoryCreateUser")
	defer span.End()

	sql, args, err := r.Builder.
		Insert("users").
		Columns("name, username, email, password_hash").
		Values(a.Name, a.Username, a.Email, a.Password).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, a entity.User) (*entity.User, error) {

	// Tracer
	tracerName := os.Getenv("OTEL_SERVICE_NAME")
	_, span := otel.GetTracerProvider().Tracer(tracerName).Start(ctx, "RepositoryGetUser")
	defer span.End()

	query, args, err := r.Builder.
		Select("*").
		From("users").
		Where("username IN (?) AND email IN (?)", a.Username, a.Email).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserRepo - Get - r.Builder: %w", err)
	}

	var id int
	var name string
	var username string
	var email string
	var password string

	row := r.Pool.QueryRow(ctx, query, args...)

	err = row.Scan(&id, &name, &username, &email, &password)
	if err != nil {
		return nil, fmt.Errorf("User not registered")
	}

	// Check password hash
	if password != a.Password {
		return nil, fmt.Errorf("Invalid credentials")
	}

	user := &entity.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	return user, nil

}
