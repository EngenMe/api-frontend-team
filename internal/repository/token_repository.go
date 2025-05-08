package repository

import (
	"github.com/EngenMe/api-frontend-team/internal/model"
)

type TokenRepository interface {
	CreateToken(token *model.Token) error
	UpdateTokenByuserId(userId string, token string) error
	FindTokenByUserId(userId string) (*model.Token, error)
}
