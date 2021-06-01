package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Credentials      Credentials `json:"credentials" gorm:"foreignKey:ID"`
	APIToken         string      `json:"apiToken" gorm:"not null"`
	IsDeleted        bool        `json:"isDeleted" gorm:"not null"`
	Roles            []Role      `json:"roles" gorm:"many2many:user_roles;"`
	Email            string      `json:"email" gorm:"not null;unique"`
	ValidationUid    string      `json:"validationUid"`
	ValidationExpire time.Time   `json:"validationExpire"`
}
