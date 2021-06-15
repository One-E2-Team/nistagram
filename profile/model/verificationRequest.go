package model

import "gorm.io/gorm"

type VerificationRequest struct {
	gorm.Model
	ProfileID          uint
	Name         	   string             `json:"name"`
	Surname      	   string             `json:"surname"`
	VerificationStatus VerificationStatus `json:"verificationStatus"`
	ImagePath          string             `json:"imagePath"`
	CategoryID		   uint
	Category 		   Category 		  `json:"category"`
}
