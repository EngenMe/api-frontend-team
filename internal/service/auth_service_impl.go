package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/pkg/jwt"
	"github.com/EngenMe/api-frontend-team/pkg/utils"
)

type authService struct {
	repo      repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(repo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{repo: repo, tokenRepo: tokenRepo}
}

func (s *authService) Login(dtoUser *dto.LoginRequest) (string, string, error) {

	// Check if user exists
	user, err := s.repo.FindByEmail(dtoUser.Email)
	if err != nil {
		return "", "", err
	}

	// Check password
	if !utils.CheckPasswordHash(dtoUser.Password, user.Password) {
		return "", "", err
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)

	// Generate JWT token
	idToken, err := jwt.GenerateToken(userIDStr, user.Email)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshIdToken, err := jwt.GenerateRefreshToken(userIDStr)
	if err != nil {
		return "", "", err
	}

	tokenModel, err := s.tokenRepo.FindTokenByUserId(userIDStr)
	if err != nil && tokenModel == nil {
		fmt.Println("Creating new token record")
		tokenModel = &model.Token{
			UserID:       userIDStr,
			RefreshToken: refreshIdToken,
		}
		err = s.tokenRepo.CreateToken(tokenModel)
		if err != nil {
			return "", "", err
		}
	} else if err != nil {
		return "", "", err
	} else {
		fmt.Println("Updating existing token record")
		err = s.tokenRepo.UpdateTokenByuserId(userIDStr, refreshIdToken)
		if err != nil {
			return "", "", err
		}
	}

	return idToken, refreshIdToken, nil
}

func (s *authService) Register(dto *dto.RegisterRequest) error {
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

func (s *authService) RefreshToken(userID string, refreshToken string) (string, string, error) {
	// Check if the refresh token is valid
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims["user_id"] != userID {
		return "", "", errors.New("invalid token")
	}

	user, _ := s.repo.GetById(userID)
	if user == nil {
		return "", "", errors.New("user not found")
	}

	// Generate a new JWT token
	newIdToken, err := jwt.GenerateToken(userID, user.Email)
	if err != nil {
		return "", "", err
	}
	// update refresh token
	newRefreshToken, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	err = s.tokenRepo.UpdateTokenByuserId(userID, newRefreshToken)
	if err != nil {
		return "", "", err
	}
	return newIdToken, newRefreshToken, nil
}
