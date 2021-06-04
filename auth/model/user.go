package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ProfileId        uint      `json:"profileId" gorm:"unique;not null"`
	Password         string    `json:"password" gorm:"not null"`
	APIToken         string    `json:"apiToken"`
	IsDeleted        bool      `json:"isDeleted" gorm:"not null"`
	IsValidated      bool      `json:"isValidated" gorm:"not null"`
	Roles            []Role    `json:"roles" gorm:"many2many:user_roles;"`
	Email            string    `json:"email" gorm:"not null;unique"`
	ValidationUid    string    `json:"validationUid"`
	ValidationExpire time.Time `json:"validationExpire"`
}
