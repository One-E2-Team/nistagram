package model

import (
	"gorm.io/gorm"
	"time"
)

type Timestamp struct {
	gorm.Model
	CampaignParametersID uint `json:"campaignParametersID"`
	Time time.Time `json:"timestamp"`
}