package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name       string      `json:"name" gorm:"unique;not null"`
	Privileges []Privilege `json:"privileges" gorm:"many2many:role_privileges;"`
}
