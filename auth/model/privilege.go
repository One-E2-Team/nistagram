package model

import "gorm.io/gorm"

type Privilege struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}
