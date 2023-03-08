package usecase

import (
	"context"

	"github.com/okassov/pet-auth/internal/entity"
)

type (
	Auth interface {
		SignUp(context.Context, entity.User) error
		SignIn(context.Context, entity.User) (map[string]string, error)
		ValidateToken(ctx context.Context, accessToken string) (UserClaim, error)
	}

	UserRepo interface {
		CreateUser(context.Context, entity.User) error
		GetUser(context.Context, entity.User) (*entity.User, error)
	}
)
