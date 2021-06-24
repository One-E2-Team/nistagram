package repository

import "gorm.io/gorm"

type CampaignRepository struct {
	Database *gorm.DB
}
