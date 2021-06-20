package model

import "gorm.io/gorm"

type AgentRequest struct {
	gorm.Model
	ProfileId uint `json:"profileId" gorm:"unique;not null"`
}
