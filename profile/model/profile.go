package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	ProfileSettings ProfileSettings `json:"profileSettings" gorm:"foreignKey:ID"`
	PersonalData PersonalData `json:"personalData" gorm:"foreignKey:ID"`
	Biography string `json:"biography"`
	Website string `json:"website"`
	Type ProfileType `json:"type"`
}
