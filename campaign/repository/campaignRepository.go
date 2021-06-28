package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/campaign/model"
	"time"
)

type CampaignRepository struct {
	Database *gorm.DB
}


func (repo *CampaignRepository) CreateCampaign(campaign model.Campaign) (model.Campaign,error) {
	result := repo.Database.Create(&campaign)
	if result.RowsAffected == 0 {
		return campaign, fmt.Errorf("User not created")
	}
	fmt.Println("User Created")
	return campaign, nil
}

func (repo *CampaignRepository) UpdateCampaignParameters(campaignParameters model.CampaignParameters)  error {

	var oldValue model.CampaignParameters
	tomorrow := time.Now().Add(24 * time.Hour)
	tx := repo.Database.Begin()
	result := tx.Table("campaign_parameters").Exec("UPDATE campaign_parameters SET end = ? WHERE id IN" +
		"(SELECT searched.id FROM (SELECT * FROM campaign_parameters cp WHERE cp.campaign_id = ? AND cp.start < ? AND cp.deleted_at IS NULL " +
		"ORDER BY cp.start DESC LIMIT 1) searched)",tomorrow, campaignParameters.CampaignID, tomorrow).Scan(&oldValue)
	if result.Error != nil {
		return result.Error
	}

	newCampParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            tomorrow,
		End:              campaignParameters.End,
		CampaignID:       campaignParameters.CampaignID,
		Interests:        campaignParameters.Interests,
		CampaignRequests: campaignParameters.CampaignRequests,
		Timestamps:       campaignParameters.Timestamps,
	}


	if err := tx.Table("campaign_parameters").Create(&newCampParams).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CampaignRepository) DeleteCampaign(campaignID uint) error{
	if err := repo.Database.Delete(&model.Campaign{},campaignID).Error; err!=nil{
		return err
	}
	return nil
}

func (repo *CampaignRepository) GetInterests(interests []string) []model.Interest {
	var ret []model.Interest

	if err := repo.Database.Table("interests").Find(&ret,"name IN ? ", interests).Error ; err != nil {
		return make([]model.Interest,0)
	}
	return ret
}

func (repo *CampaignRepository) GetParametersByCampaignId(campaignId uint) (model.CampaignParameters, error) {
	var ret model.CampaignParameters

	err := repo.Database.Table("campaign_parameters").Preload("Interests").
		Find(&ret).Where("campaign_id = ?", campaignId).
		Where("start <= ?", time.Now()).Where("end >= ?", time.Now()).Error

	return ret, err
}
