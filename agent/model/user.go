package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID               uint 			      `json:"id" gorm:"primaryKey"`
	Email            string               `json:"email" gorm:"not null;unique"`
	Password         string               `json:"password" gorm:"not null"`
	APIToken         string               `json:"apiToken"`
	Address          string               `json:"address"`
	IsValidated      bool                 `json:"isValidated" gorm:"not null"`
	Roles            []Role               `json:"roles" gorm:"many2many:user_roles;"`
	ValidationUid    string               `json:"validationUid"`
	ValidationExpire time.Time            `json:"validationExpire"`
}

type Role struct {
	gorm.Model
	Name       string      `json:"name" gorm:"unique;not null"`
	Privileges []Privilege `json:"privileges" gorm:"many2many:role_privileges;"`
}

type Privilege struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}
