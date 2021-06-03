package dto

type UpdateUserDTO struct {
	ProfileId string `json:"profileId" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"not null;unique"`
}
