package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	ProfileSettings ProfileSettings `json:"profileSettings"`
	PersonalData PersonalData `json:"personalData"`
	Biography string `json:"biography"`
	Website string `json:"website"`
	Type ProfileType `json:"type"`
}
