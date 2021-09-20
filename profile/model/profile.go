package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username        string          `json:"username" gorm:"unique;not null"`
	Email           string          `json:"email" gorm:"unique; not null"`
	ProfileSettings ProfileSettings `json:"profileSettings" gorm:"embedded"`
	PersonalData    PersonalData    `json:"personalData" gorm:"embedded"`
	Biography       string          `json:"biography"`
	Website         string          `json:"website"`
	IsVerified      bool            `json:"isVerified"`
	InterestedIn []Interest `json:"interestedIn" gorm:"many2many:person_interests;"`
}

func (profile *Profile) AddItem(item Interest) {
	profile.InterestedIn = append(profile.InterestedIn, item)
}
