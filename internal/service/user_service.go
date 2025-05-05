package service

import "github.com/EngenMe/api-frontend-team/internal/model"

type UserService interface {
	GetUserByEmail(email string) (*model.User, error)
}
