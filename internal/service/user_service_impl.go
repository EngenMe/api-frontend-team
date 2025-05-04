package service

import (
	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(email, password string) error {
	user := &model.User{
		Email:    email,
		Password: password, // In production, hash the password
	}
	return s.repo.Create(user)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
}
