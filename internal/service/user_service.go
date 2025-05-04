package service

import "github.com/EngenMe/api-frontend-team/internal/model"

type UserService interface {
	Register(email, password string) error
	GetUserByEmail(email string) (*model.User, error)
}
