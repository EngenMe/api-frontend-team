package repository

import (
	"github.com/EngenMe/api-frontend-team/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	GetById(id string) (*model.User, error)
	UpdateUser(id string, user *model.User) (*model.User, error)
	DeleteUser(id string) error
}
