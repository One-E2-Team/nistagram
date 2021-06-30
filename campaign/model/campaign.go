package model

import (
	"gorm.io/gorm"
	"time"
)

type Campaign struct {
	gorm.Model
	PostID             string               `json:"postId" gorm:"not null"`
	AgentID            uint                 `json:"agentId" gorm:"not null"`
	CampaignType       CampaignType         `json:"campaignType" gorm:"not null"`
	Start              time.Time            `json:"start" gorm:"not null"`
	CampaignParameters []CampaignParameters `json:"campaignParameters" gorm:"constraint:OnDelete:CASCADE"`
}
