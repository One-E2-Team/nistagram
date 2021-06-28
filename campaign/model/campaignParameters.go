package model

import (
	"gorm.io/gorm"
	"time"
)

type CampaignParameters struct {
	gorm.Model
	Start            time.Time         `json:"start" gorm:"not null"`
	End              time.Time         `json:"end" gorm:"not null"`
	CampaignID       uint              `json:"campaignId" gorm:"not null"`
	Interests        []Interest        `json:"interests" gorm:"many2many:parameters_interests;constraint:OnDelete:CASCADE"`
	CampaignRequests []CampaignRequest `json:"campaignRequests"  gorm:"foreignKey:campaign_parameters_id;constraint:OnDelete:CASCADE"`
	Timestamps  	 []Timestamp	   `json:"timestamps" gorm:"not null"`
}


