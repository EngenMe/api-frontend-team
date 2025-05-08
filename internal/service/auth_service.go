package service

import "github.com/EngenMe/api-frontend-team/internal/dto"

type AuthService interface {
	Register(request *dto.RegisterRequest) error
	Login(request *dto.LoginRequest) (string, string, error)
	RefreshToken(userID string, refreshToken string) (string, string, error)
}
