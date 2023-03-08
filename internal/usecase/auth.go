package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/okassov/pet-auth/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	jwt "github.com/dgrijalva/jwt-go/v4"
)

type AccessAuthClaims struct {
	jwt.StandardClaims
	User UserClaim `json:"user"`
}

type RefreshAuthClaims struct {
	jwt.StandardClaims
}

type UserClaim struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AuthUseCase struct {
	repo           UserRepo
	signingKey     []byte
	expireDuration time.Duration
	tracer         trace.TracerProvider
	tracerName     string
}

func New(r UserRepo, signingKey []byte, tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		repo:           r,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, u entity.User) error {

	// Tracer
	tracerName := os.Getenv("OTEL_SERVICE_NAME")
	newCtx, span := otel.GetTracerProvider().Tracer(tracerName).Start(ctx, "AuthUseCaseSignUp")

	err := uc.repo.CreateUser(newCtx, entity.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: GeneratePasswordHash(u.Password),
	})
	if err != nil {
		span.End()
		return err
	}

	span.End()
	return nil
}

func (uc *AuthUseCase) SignIn(ctx context.Context, u entity.User) (map[string]string, error) {

	// Tracer
	tracerName := os.Getenv("OTEL_SERVICE_NAME")
	newCtx, span := otel.GetTracerProvider().Tracer(tracerName).Start(ctx, "AuthUseCaseSignIn")

	user, err := uc.repo.GetUser(newCtx, entity.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: GeneratePasswordHash(u.Password),
	})

	if err != nil {
		span.End()
		return nil, err
	}

	// JWT Access Token
	accessTokenClaims := AccessAuthClaims{
		User: UserClaim{
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(uc.expireDuration)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	at, err := accessToken.SignedString(uc.signingKey)
	if err != nil {
		span.End()
		return nil, err
	}

	// JWT Refresh Token
	refreshTokenClaims := RefreshAuthClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(uc.expireDuration)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	rt, err := refreshToken.SignedString(uc.signingKey)
	if err != nil {
		span.End()
		return nil, err
	}

	span.End()
	return map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	}, nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, accessToken string) (UserClaim, error) {

	token, err := jwt.ParseWithClaims(accessToken, &AccessAuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return uc.signingKey, nil
	})

	if err != nil {
		// return nil, err
		return UserClaim{}, err
	}

	if claims, ok := token.Claims.(*AccessAuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	// return nil, err
	return UserClaim{}, nil
}
