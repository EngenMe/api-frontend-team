package service

import (
	"errors"
	"strconv"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/pkg/jwt"
	"github.com/EngenMe/api-frontend-team/pkg/utils"
)

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Login(dto *dto.UserDTO) (string, error) {
	// Check if user exists
	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		return "", err
	}

	// Check password
	if !utils.CheckPasswordHash(dto.Password, user.Password) {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Register(dto *dto.UserDTO) error {
	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(dto.Email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}
	hPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	if !isEmailValid(dto.Email) {
		return errors.New("invalid email format")
	}

	user := &model.User{
		Email:    dto.Email,
		Password: hPassword, // In production, hash the password
	}
	return s.repo.Create(user)
}
