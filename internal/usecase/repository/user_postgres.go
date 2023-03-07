package repository

import (
	"context"
	"fmt"

	"github.com/okassov/pet-auth/internal/entity"
	"github.com/okassov/pet-auth/pkg/postgres"
)

type User struct {
	// ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"password"`
	Password string `json:"password"`
}

type UserRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

//
func (r *UserRepo) CreateUser(ctx context.Context, a entity.User) error {

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
		// return fmt.Errorf("UserRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, a entity.User) (*entity.User, error) {

	query, args, err := r.Builder.
		Select("*").
		From("users").
		Where("username IN (?) AND password_hash IN (?)", a.Username, a.Password).
		ToSql()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("UserRepo - Get - r.Builder: %w", err)
	}

	var id int
	var name string
	var username string
	var email string
	var password_hash string

	row := r.Pool.QueryRow(ctx, query, args...) //.Scan(&id)
	// values := row.Values()
	err = row.Scan(&id, &name, &username, &email, &password_hash)
	if err != nil {
		return nil, fmt.Errorf("User not registered")
	}

	return toModel(&User{Name: name, Username: username, Email: email, Password: password_hash}), nil
}

func toModel(u *User) *entity.User {
	return &entity.User{
		// ID:       u.ID.Hex(),
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
