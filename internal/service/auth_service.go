package service

import "github.com/EngenMe/api-frontend-team/internal/dto"

type AuthService interface {
	Register(request *dto.RegisterRequest) (dto.AuthUserResponse, error)
	Login(request *dto.LoginRequest) (dto.AuthUserResponse, error)
	RefreshToken(userID string, refreshToken string) (dto.RefreshTokenResonse, error)
}
