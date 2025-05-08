package service

import "github.com/EngenMe/api-frontend-team/internal/dto"

type AuthService interface {
	Register(request *dto.RegisterRequest) error
	Login(request *dto.LoginRequest) (dto.RefreshTokenResonse, error)
	RefreshToken(userID string, refreshToken string) (dto.RefreshTokenResonse, error)
}
