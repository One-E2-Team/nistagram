package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/campaign/model"
	"nistagram/util"
	"strings"
	"time"
)

type CampaignRepository struct {
	Database *gorm.DB
}


func (repo *CampaignRepository) CreateCampaign(campaign model.Campaign) (model.Campaign,error) {
	result := repo.Database.Create(&campaign)
	if result.RowsAffected == 0 {
		return campaign, fmt.Errorf("Campaign not created")
	}
	fmt.Println("Campaign Created")
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
	now := time.Now()
	//Check if campaign exists
	if err := repo.checkIfCampaignExists(campaignID); err != nil{
		return err
	}

	tx := repo.Database.Begin()

	//check if campaign have params in past
	if result := tx.Table("campaign_parameters").Find(&model.CampaignParameters{}, "start < ? AND campaign_id = ?",now,campaignID); result.Error != nil{
		return result.Error
	}else if result.RowsAffected == 0 {
		//if there is no params delete campaign
		if err := repo.forceDeleteCampaing(campaignID,tx); err != nil{
			return err
		}
		tx.Commit()
		return nil
	}

	//delete all future campaign params
	var deletedCampaignParameters []model.CampaignParameters
	if result := tx.Find(&deletedCampaignParameters,"start > ? AND campaign_id = ?", time.Now(), campaignID); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	ids := make([]uint , 0)
	for _, value := range deletedCampaignParameters{
		ids = append(ids,value.ID)
	}
	if err:=beforeDeleteCampaignParameters(ids,tx); err != nil{
		tx.Rollback()
		return err
	}

	if result := tx.Unscoped().Delete(&[]model.CampaignParameters{},"id IN ?",ids); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//Set active CampaignParameter to end now
	result := tx.Table("campaign_parameters").Exec("UPDATE campaign_parameters SET end = ? WHERE id IN"+
		"(SELECT searched.id FROM (SELECT * FROM campaign_parameters cp WHERE cp.campaign_id = ? AND cp.end > ? AND cp.deleted_at IS NULL "+
		"ORDER BY cp.start DESC LIMIT 1) searched)", time.Now(), campaignID, time.Now()).Scan(&model.CampaignParameters{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (repo *CampaignRepository) GetMyCampaigns(agentID uint) ([]model.Campaign, error) {
	var ret []model.Campaign

	if err := repo.Database.Table("campaigns").Find(&ret,"agent_id = ? ", agentID).Error ; err != nil {
		return make([]model.Campaign,0), err
	}
	return ret, nil
}

func (repo *CampaignRepository) GetAllInterests() ([]string, error) {
	var interests []string
	result := repo.Database.Table("interests").Select("name").Find(&interests, "name LIKE ?", "%%")
	return interests, result.Error
}

func (repo *CampaignRepository) checkIfCampaignExists(campaignID uint) error {
	if result := repo.Database.Find(&model.Campaign{},"id = ?", campaignID); result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
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

func (repo *CampaignRepository) forceDeleteCampaing(id uint, tx *gorm.DB) error {
	if err := deleteCampaignParameters(id,tx); err != nil {
		return err
	}
	return tx.Unscoped().Delete(&model.Campaign{},"id = ?",id).Error
}

func  deleteCampaignParameters(campaignId uint, tx *gorm.DB) error {
	var ret []model.CampaignParameters
	if err := tx.Table("campaign_parameters").Find(&ret,"campaign_id = ? ", campaignId).Error ; err != nil {
		return err
	}
	ids := make([]uint,0)
	for _, value := range ret {
		ids = append(ids,value.ID)
	}
	if err := beforeDeleteCampaignParameters(ids,tx); err != nil {
		return err
	}
	return tx.Unscoped().Delete(&model.CampaignParameters{},"campaign_id = ?",campaignId).Error
}

func  beforeDeleteCampaignParameters(campaignParametersId []uint,tx *gorm.DB) error {
	if err:= tx.Unscoped().Delete(&model.CampaignRequest{},"campaign_parameters_id IN ?",campaignParametersId).Error; err != nil{
		return err
	}
	if err:= tx.Unscoped().Delete(&model.Timestamp{},"campaign_parameters_id IN ?",campaignParametersId).Error; err != nil{
		return err
	}

	return nil
}


func (repo *CampaignRepository) GetParametersByCampaignId(campaignId uint) (model.CampaignParameters, error) {
	var ret model.CampaignParameters

	err := repo.Database.Table("campaign_parameters").Preload("Interests").
		Find(&ret).Where("campaign_id = ?", campaignId).
		Where("start <= ?", time.Now()).Where("end >= ?", time.Now()).Error

	return ret, err
}

func (repo *CampaignRepository) GetCampaignById(campaignId uint) (model.Campaign, error) {
	var ret model.Campaign

	err := repo.Database.Preload("CampaignParameters.Interests").
		Preload("CampaignParameters.CampaignRequests").
		Preload("CampaignParameters.Timestamps").
		Find(&ret).Where("id = ?", campaignId).Error

	return ret, err
}

func (repo *CampaignRepository) GetLastActiveParametersForCampaign(campaignId uint) (model.CampaignParameters, error) {
	var ret model.CampaignParameters
	result := repo.Database.Preload("Interests").Preload("CampaignRequests").Preload("Timestamps").
		Raw("select * from campaign_parameters p where p.campaign_id = ? AND p.end > NOW() AND deleted_at IS NULL AND p.start < NOW()", campaignId).First(&ret)
	if result.Error != nil {
		return model.CampaignParameters{},result.Error
	}else if result.RowsAffected == 0 {
		return model.CampaignParameters{}, gorm.ErrRecordNotFound
	}
	return ret, nil
}

func (repo *CampaignRepository) GetAllActiveParameters() ([]model.CampaignParameters, error) {
	var ret []model.CampaignParameters
	result := repo.Database.Preload("Interests").Preload("CampaignRequests").Preload("Timestamps").
		Where("end > ? AND deleted_at IS NULL AND start < ? ", time.Now(), time.Now()).Find(&ret)

	if result.Error != nil {
		return make([]model.CampaignParameters, 0), result.Error
	}else if result.RowsAffected == 0 {
		return make([]model.CampaignParameters, 0), gorm.ErrRecordNotFound
	}

	return ret, nil
}

func (repo *CampaignRepository) GetPostIDsFromCampaignIDs(campaignIDs []uint) ([]string, error) {
	var ret []string
	var builder strings.Builder
	// https://stackoverflow.com/questions/39348610/select-query-with-in-clause-having-duplicate-values-in-in-clause
	builder.WriteString("select c.post_id from campaigns c join (")
	length := len(campaignIDs)
	for i, campaignID := range campaignIDs {
		builder.WriteString("select " + util.Uint2String(campaignID) + " ")
		if i == 0 {
			builder.WriteString("as id ")
		}
		if i != length - 1 {
			builder.WriteString("union all ")
		}
	}
	builder.WriteString(") matches using (id)")
	fmt.Println("query: ", builder.String())
	result := repo.Database.Raw(builder.String()).Scan(&ret)
	if result.Error != nil {
		return make([]string, 0), result.Error
	}else if result.RowsAffected == 0 {
		return make([]string, 0), gorm.ErrRecordNotFound
	}
	return ret, nil
}

func (repo *CampaignRepository) UpdateCampaignRequest(id string, status model.RequestStatus) error {
	if res := repo.Database.Model(&model.CampaignRequest{}).Where("id = ?", id).Update("request_status", status); res.Error != nil{
		return res.Error
	}else if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *CampaignRepository) GetDistinctCampaignParamsIdForProfileId(id int) ([]int,error) {
	var ret []int
	if result := repo.Database.Raw("select distinct campaign_parameters_id " +
		"FROM campaign_requests " +
		"WHERE influencer_id = ? " +
		"AND request_status = ?", id,model.SENT).Scan(&ret); result.Error != nil{
		return make([]int,0),result.Error
	}
	return ret,nil
}

func (repo *CampaignRepository) GetActiveCampaignIdsForCampaignParamsIds(campaignParamsIds []int) ([]int,error) {
	var ret []int
	result := repo.Database.Raw("select distinct campaign_id " +
		"FROM campaign_parameters " +
		"WHERE id IN (?) AND " +
		"end > ? AND deleted_at IS NULL AND start < ? ", campaignParamsIds, time.Now(), time.Now()).Scan(&ret)
	if result.Error != nil {
		return make([]int,0) , nil
	}
	return ret,nil
}

func (repo *CampaignRepository) GetCampaignRequestInfluencerId(requestId uint) uint {
	var ret uint
	res := repo.Database.Raw("select influencer_id from campaign_requests" +
		"where id = ? ", requestId).Scan(&ret)
	if res.RowsAffected == 0 {
		return 0
	}
	return ret
}
