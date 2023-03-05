package usecase

import "context"

type AuthUseCase struct {
}

func New() *AuthUseCase {
	return &AuthUseCase{}
}

func (uc *AuthUseCase) GenerateToken(ctx context.Context) {
	return
}
