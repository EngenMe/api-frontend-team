package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/EngenMe/api-frontend-team/internal/dto"
	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/pkg/utils"
)

type userService struct {
	repo repository.UserRepository
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// func (s *userService) GetUserByEmail(email string) (dto.GetUserResponse, error) {
// 	user, err := s.repo.FindByEmail(email)
// 	if err != nil {
// 		return dto.GetUserResponse{}, err
// 	}
// 	if user == nil {
// 		return dto.GetUserResponse{}, errors.New("user not found")
// 	}
// 	return dto.GetUserResponse{
// 		Email: user.Email,
// 	}, nil
// }

func (s *userService) GetUserById(id string) (dto.GetUserResponse, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return dto.GetUserResponse{}, err
	}
	if user == nil {
		return dto.GetUserResponse{}, errors.New("user not found")
	}
	return dto.GetUserResponse{
		Id:    fmt.Sprintf("%d", user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *userService) DeleteUser(id string) error {
	_, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteUser(id)
}

func (s *userService) UpdateUser(id string, userDTO dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	if !isEmailValid(userDTO.Email) {
		return dto.UpdateUserResponse{}, errors.New("invalid email format")
	}
	if userDTO.Password != "" {
		hPassword, err := utils.HashPassword(userDTO.Password)
		if err != nil {
			return dto.UpdateUserResponse{}, err
		}
		userDTO.Password = hPassword
	}
	user := &model.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	}

	UpdateUser, err := s.repo.UpdateUser(id, user)
	if err != nil {
		return dto.UpdateUserResponse{}, err
	}

	return dto.UpdateUserResponse{
		Id:    id,
		Name:  UpdateUser.Name,
		Email: UpdateUser.Email}, nil
}
