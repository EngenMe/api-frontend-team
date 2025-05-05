package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `json:"email" gorm:"unique;not null"`
	//TODO: A mistake here to json:"password" it should be something else! why...?
	Password string `json:"password" gorm:"not null"`
}
