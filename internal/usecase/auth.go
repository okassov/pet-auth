package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/okassov/pet-auth/internal/entity"

	jwt "github.com/dgrijalva/jwt-go/v4"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *entity.User `json:"user"`
}

type AuthUseCase struct {
	repo           UserRepo
	signingKey     []byte
	expireDuration time.Duration
}

func New(r UserRepo, signingKey []byte, tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		repo:           r,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, u entity.User) error {

	hashPassword := GeneratePasswordHash(u.Password)

	user := entity.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: hashPassword,
	}

	err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) SignIn(ctx context.Context, u entity.User) (string, error) {

	hashPassword := GeneratePasswordHash(u.Password)

	user := entity.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: hashPassword,
	}

	getUser, err := uc.repo.GetUser(ctx, user)
	if err != nil {
		return "", err
	}

	claims := AuthClaims{
		User: getUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(uc.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uc.signingKey)
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, accessToken string) (*entity.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return uc.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, err
}
