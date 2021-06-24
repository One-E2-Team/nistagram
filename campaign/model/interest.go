package model

import "gorm.io/gorm"

type Interest struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}
