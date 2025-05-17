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

func (s *authService) Login(dtoUser *dto.LoginRequest) (dto.AuthUserResponse, error) {
	var authUserResponse dto.AuthUserResponse
	// Check if user exists
	user, err := s.repo.FindByEmail(dtoUser.Email)
	if err != nil {
		return authUserResponse, err
	}

	// Check password
	if !utils.CheckPasswordHash(dtoUser.Password, user.Password) {
		return authUserResponse, err
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)

	// Generate JWT token

	refresh_token_res, err := refreshTokens(userIDStr, user.Email)
	if err != nil {
		return authUserResponse, err
	}

	tokenModel, err := s.tokenRepo.FindTokenByUserId(userIDStr)
	if err != nil && tokenModel == nil {
		fmt.Println("Creating new token record")
		tokenModel = &model.Token{
			UserID:       userIDStr,
			RefreshToken: refresh_token_res.Refresh.Token,
		}
		err = s.tokenRepo.CreateToken(tokenModel)
		if err != nil {
			return authUserResponse, err
		}
	} else if err != nil {
		return authUserResponse, err
	} else {
		fmt.Println("Updating existing token record")
		err = s.tokenRepo.UpdateTokenByuserId(userIDStr, refresh_token_res.Refresh.Token)
		if err != nil {
			return authUserResponse, err
		}
	}

	authUserResponse.Tokens = refresh_token_res
	authUserResponse.User = dto.User{
		Id:    userIDStr,
		Name:  user.Name,
		Email: user.Email,
	}
	return authUserResponse, nil
}

func (s *authService) Register(userDto *dto.RegisterRequest) (dto.AuthUserResponse, error) {

	var authUserResponse dto.AuthUserResponse

	existingUser, err := s.repo.FindByEmail(userDto.Email)
	if err == nil && existingUser != nil {
		return authUserResponse, errors.New("user already exists")
	}
	hPassword, err := utils.HashPassword(userDto.Password)
	if err != nil {
		return authUserResponse, err
	}

	if !isEmailValid(userDto.Email) {
		return authUserResponse, errors.New("invalid email format")
	}

	user := &model.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: hPassword, // In production, hash the password
	}
	err = s.repo.Create(user)
	if err != nil {
		return authUserResponse, errors.New("can't create user")
	}

	createdUser, _ := s.repo.FindByEmail(user.Email)
	userIdStr := strconv.FormatUint(uint64(createdUser.ID), 10)
	tokens, _ := refreshTokens(userIdStr, createdUser.Email)

	authUserResponse.Tokens = tokens
	authUserResponse.User = dto.User{
		Id:    userIdStr,
		Name:  user.Name,
		Email: user.Email,
	}
	return authUserResponse, nil
}

func (s *authService) RefreshToken(userID string, refreshToken string) (dto.RefreshTokenResonse, error) {

	var refresh_token dto.RefreshTokenResonse

	savedToken, err := s.tokenRepo.FindTokenByUserId(userID)
	if err != nil || savedToken.RefreshToken != refreshToken {
		return refresh_token, errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.GetById(userID)
	if err != nil || user == nil {
		return refresh_token, errors.New("user not found")
	}

	refresh_token_res, err := refreshTokens(userID, user.Email)
	if err != nil {
		return refresh_token, err
	}
	// update refresh token
	err = s.tokenRepo.UpdateTokenByuserId(userID, refresh_token_res.Refresh.Token)
	if err != nil {
		return refresh_token, err
	}

	return refresh_token_res, nil
}

func refreshTokens(userId string, userEmail string) (dto.RefreshTokenResonse, error) {
	var refresh_token dto.RefreshTokenResonse

	// Generate a new JWT token
	idTokenRes, idExpRes, err := jwt.GenerateToken(userId, userEmail)
	if err != nil {
		return refresh_token, err
	}
	accessPair := dto.TokenPair{
		Token:   idTokenRes,
		Expires: idExpRes,
	}

	idRefreshTokenRes, refreshExpRes, err := jwt.GenerateRefreshToken(userId)
	if err != nil {
		return refresh_token, err
	}

	refreshPair := dto.TokenPair{
		Token:   idRefreshTokenRes,
		Expires: refreshExpRes,
	}

	refresh_token.Access = accessPair
	refresh_token.Refresh = refreshPair

	return refresh_token, nil
}
