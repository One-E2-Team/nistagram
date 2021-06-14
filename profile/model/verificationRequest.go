package model

import "gorm.io/gorm"

type VerificationRequest struct {
	gorm.Model
	Profile            Profile            `gorm:"foreignKey:ID"`
	VerificationStatus VerificationStatus `json:"verificationStatus"`
	Category           Category           `json:"category" gorm:"foreignKey:ID"`
	ImagePath          string             `json:"imagePath"`
}
