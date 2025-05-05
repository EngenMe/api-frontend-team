package service

import (
	"regexp"

	"github.com/EngenMe/api-frontend-team/internal/model"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"golang.org/x/crypto/bcrypt"

	"errors"
)

// TODO: Move register and login to auth_service, create interface for that DI principle
type userService struct {
	repo repository.UserRepository
}

//TODO: Refactoring: Move hash password and check password hash to pkg/utils/password_helper.go or something like that

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(email, password string) error {
	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("user already exists")
	}

	hPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	if !isEmailValid(email) {
		return errors.New("invalid email format")
	}

	user := &model.User{
		Email:    email,
		Password: hPassword, // In production, hash the password
	}
	return s.repo.Create(user)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
}
