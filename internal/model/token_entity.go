package model

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	UserID       string `json:"userid" gorm:"unique;not null"`
	RefreshToken string `json:"RefreshToken" gorm:"not null"`
}
