package service

import "github.com/EngenMe/api-frontend-team/internal/dto"

type AuthService interface {
	Register(*dto.UserDTO) error
	Login(*dto.UserDTO) (string, error)
}
