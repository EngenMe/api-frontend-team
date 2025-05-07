package repository

import (
	"github.com/EngenMe/api-frontend-team/internal/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetById(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(id string, user *model.User) (*model.User, error) {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	err := r.db.Unscoped().Delete(&model.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
