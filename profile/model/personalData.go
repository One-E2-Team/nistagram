package model

import "gorm.io/gorm"

type PersonalData struct {
	gorm.Model
	Name         string     `json:"name"`
	Surname      string     `json:"surname"`
	Telephone    string     `json:"telephone"`
	Gender       string     `json:"gender"`
	BirthDate    string     `json:"birthDate"`
	InterestedIn []Interest `json:"interestedIn" gorm:"many2many:person_interests;"`
	ProfileID    uint
}

func (personalData *PersonalData) AddItem(item Interest) {
	personalData.InterestedIn = append(personalData.InterestedIn, item)
}
