package service

import "github.com/EngenMe/api-frontend-team/internal/model"

type UserService interface {
	// TODO: Move register and login to auth_service, create interface for that DI principle
	Register(email, password string) error
	GetUserByEmail(email string) (*model.User, error)
	Login(email, password string) (string, error)
}
