package service

import (
	"errors"
	"strconv"

	"github.com/EngenMe/api-frontend-team/pkg/jwt"
)

func (s *userService) Login(email, password string) (string, error) {
	// Check if user exists
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	// Check password
	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
