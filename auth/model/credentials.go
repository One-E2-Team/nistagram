package model

import "gorm.io/gorm"

type Credentials struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Salt     string `json:"salt" gorm:"not null"`
}
