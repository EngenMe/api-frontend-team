package repository

import "github.com/EngenMe/api-frontend-team/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}
