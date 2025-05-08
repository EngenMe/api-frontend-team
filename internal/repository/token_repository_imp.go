package repository

import (
	"github.com/EngenMe/api-frontend-team/internal/model"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) TokenRepository {
	db.AutoMigrate(&model.Token{})
	if err := db.AutoMigrate(&model.Token{}); err != nil {
		panic("failed to migrate token table")
	}
	return &tokenRepository{db}
}

func (t *tokenRepository) CreateToken(token *model.Token) error {
	return t.db.Create(token).Error
}

func (t *tokenRepository) UpdateTokenByuserId(userId string, updatedtoken string) error {
	var token *model.Token
	err := t.db.Where("user_id = ?", userId).Find(&token).Error
	if err != nil {
		return err
	}

	err = t.db.Model(&token).Update("refresh_token", updatedtoken).Where("user_id = ?", userId).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *tokenRepository) FindTokenByUserId(userID string) (*model.Token, error) {
	var token model.Token
	err := t.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}
