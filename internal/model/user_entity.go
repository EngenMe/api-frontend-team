package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name       string `json:"name"  gorm:"not null"`
	Email      string `json:"email" gorm:"unique;not null"`
	Password   string `json:"-" gorm:"default:null"`
	Provider   string `json:"provider" gorm:"not null;default:'local'"`
	ProviderID string `json:"provider_id" gorm:"index default:null"`
}
