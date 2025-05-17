package service

import (
	"github.com/EngenMe/api-frontend-team/internal/dto"
)

type UserService interface {
	// GetUserByEmail(email string) (dto.GetUserResponse, error)
	GetUserById(id string) (dto.GetUserResponse, error)
	DeleteUser(id string) error
	UpdateUser(id string, user dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
}
