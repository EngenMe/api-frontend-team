package service

import (
	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/markbates/goth"
)

type AuthService interface {
	Register(request *dto.RegisterRequest) (dto.AuthUserResponse, error)
	Login(request *dto.LoginRequest) (dto.AuthUserResponse, error)
	RefreshToken(userID string, refreshToken string) (dto.RefreshTokenResonse, error)
	ProviderLogin(gothUser goth.User) (dto.AuthUserResponse, error)
}
