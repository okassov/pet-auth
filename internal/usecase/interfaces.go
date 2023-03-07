package usecase

import (
	"context"

	"github.com/okassov/pet-auth/internal/entity"
)

type (
	Auth interface {
		SignUp(context.Context, entity.User) error
		SignIn(context.Context, entity.User) (string, error)
		ValidateToken(ctx context.Context, accessToken string) (*entity.User, error)
	}

	UserRepo interface {
		CreateUser(context.Context, entity.User) error
		GetUser(context.Context, entity.User) (*entity.User, error)
	}
)
