package repository

import (
	"context"
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


func (repo *CampaignRepository) CreateCampaign(ctx context.Context, campaign model.Campaign) (model.Campaign,error) {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateCampaign-repository")
	defer util.Tracer.FinishSpan(span)

	result := repo.Database.Create(&campaign)
	if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("campaign not created"))
		return campaign, fmt.Errorf("campaign not created")
	}
	fmt.Println("Campaign Created")
	return campaign, nil
}

func (repo *CampaignRepository) UpdateCampaignParameters(ctx context.Context, campaignParameters model.CampaignParameters)  error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateCampaignParameters-repository")
	defer util.Tracer.FinishSpan(span)

	var oldValue model.CampaignParameters
	tomorrow := time.Now().Add(24 * time.Hour)
	tx := repo.Database.Begin()
	result := tx.Table("campaign_parameters").Exec("UPDATE campaign_parameters SET end = ? WHERE id IN" +
		"(SELECT searched.id FROM (SELECT * FROM campaign_parameters cp WHERE cp.campaign_id = ? AND cp.start < ? AND cp.deleted_at IS NULL " +
		"ORDER BY cp.start DESC LIMIT 1) searched)",tomorrow, campaignParameters.CampaignID, tomorrow).Scan(&oldValue)
	if result.Error != nil {
		util.Tracer.LogError(span, result.Error)
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
		util.Tracer.LogError(span, err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repo *CampaignRepository) DeleteCampaign(ctx context.Context, campaignID uint) error{
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteCampaign-repository")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	now := time.Now()

	if err := repo.checkIfCampaignExists(nextCtx, campaignID); err != nil{
		util.Tracer.LogError(span, err)
		return err
	}

	tx := repo.Database.Begin()

	//check if campaign have params in past
	if result := tx.Table("campaign_parameters").Find(&model.CampaignParameters{}, "start < ? AND campaign_id = ?",now,campaignID); result.Error != nil{
		util.Tracer.LogError(span, result.Error)
		return result.Error
	}else if result.RowsAffected == 0 {
		//if there is no params delete campaign
		if err := repo.forceDeleteCampaign(nextCtx, campaignID,tx); err != nil{
			util.Tracer.LogError(span, err)
			return err
		}
		tx.Commit()
		return nil
	}

	//delete all future campaign params
	var deletedCampaignParameters []model.CampaignParameters
	if result := tx.Find(&deletedCampaignParameters,"start > ? AND campaign_id = ?", time.Now(), campaignID); result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		tx.Rollback()
		return result.Error
	}
	ids := make([]uint , 0)
	for _, value := range deletedCampaignParameters{
		ids = append(ids,value.ID)
	}
	if err:=beforeDeleteCampaignParameters(nextCtx, ids,tx); err != nil{
		util.Tracer.LogError(span, err)
		tx.Rollback()
		return err
	}

	if result := tx.Unscoped().Delete(&[]model.CampaignParameters{},"id IN ?",ids); result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		tx.Rollback()
		return result.Error
	}

	//Set active CampaignParameter to end now
	result := tx.Table("campaign_parameters").Exec("UPDATE campaign_parameters SET end = ? WHERE id IN"+
		"(SELECT searched.id FROM (SELECT * FROM campaign_parameters cp WHERE cp.campaign_id = ? AND cp.end > ? AND cp.deleted_at IS NULL "+
		"ORDER BY cp.start DESC LIMIT 1) searched)", time.Now(), campaignID, time.Now()).Scan(&model.CampaignParameters{})
	if result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (repo *CampaignRepository) GetMyCampaigns(ctx context.Context, agentID uint) ([]model.Campaign, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyCampaigns-repository")
	defer util.Tracer.FinishSpan(span)

	var ret []model.Campaign

	if err := repo.Database.Table("campaigns").Find(&ret,"agent_id = ? ", agentID).Error ; err != nil {
		util.Tracer.LogError(span, err)
		return make([]model.Campaign,0), err
	}
	return ret, nil
}

func (repo *CampaignRepository) GetAllInterests(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllInterests-repository")
	defer util.Tracer.FinishSpan(span)

	var interests []string
	result := repo.Database.Table("interests").Select("name").Find(&interests, "name LIKE ?", "%%")
	if result.Error != nil {
		util.Tracer.LogError(span, result.Error)
	}
	return interests, result.Error
}

func (repo *CampaignRepository) checkIfCampaignExists(ctx context.Context, campaignID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "checkIfCampaignExists-repository")
	defer util.Tracer.FinishSpan(span)

	if result := repo.Database.Find(&model.Campaign{},"id = ?", campaignID); result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		return result.Error
	} else if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("campaign with id %v doesn't exist", campaignID))
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *CampaignRepository) GetInterests(ctx context.Context, interests []string) []model.Interest {
	span := util.Tracer.StartSpanFromContext(ctx, "GetInterests-repository")
	defer util.Tracer.FinishSpan(span)

	var ret []model.Interest

	if err := repo.Database.Table("interests").Find(&ret,"name IN ? ", interests).Error ; err != nil {
		util.Tracer.LogError(span, err)
		return make([]model.Interest,0)
	}
	return ret
}

func (repo *CampaignRepository) forceDeleteCampaign(ctx context.Context, id uint, tx *gorm.DB) error {
	span := util.Tracer.StartSpanFromContext(ctx, "forceDeleteCampaign-repository")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if err := deleteCampaignParameters(nextCtx, id,tx); err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return tx.Unscoped().Delete(&model.Campaign{},"id = ?",id).Error
}

func deleteCampaignParameters(ctx context.Context, campaignId uint, tx *gorm.DB) error {
	span := util.Tracer.StartSpanFromContext(ctx, "deleteCampaignParameters-repository")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var ret []model.CampaignParameters
	if err := tx.Table("campaign_parameters").Find(&ret,"campaign_id = ? ", campaignId).Error ; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	ids := make([]uint,0)
	for _, value := range ret {
		ids = append(ids,value.ID)
	}
	if err := beforeDeleteCampaignParameters(nextCtx, ids,tx); err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return tx.Unscoped().Delete(&model.CampaignParameters{},"campaign_id = ?",campaignId).Error
}

func beforeDeleteCampaignParameters(ctx context.Context, campaignParametersId []uint,tx *gorm.DB) error {
	span := util.Tracer.StartSpanFromContext(ctx, "beforeDeleteCampaignParameters-repository")
	defer util.Tracer.FinishSpan(span)

	if err:= tx.Unscoped().Delete(&model.CampaignRequest{},"campaign_parameters_id IN ?",campaignParametersId).Error; err != nil{
		util.Tracer.LogError(span, err)
		return err
	}
	if err:= tx.Unscoped().Delete(&model.Timestamp{},"campaign_parameters_id IN ?",campaignParametersId).Error; err != nil{
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *CampaignRepository) GetParametersByCampaignId(ctx context.Context, campaignId uint) (model.CampaignParameters, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetParametersByCampaignId-repository")
	defer util.Tracer.FinishSpan(span)

	var ret model.CampaignParameters

	err := repo.Database.Table("campaign_parameters").Preload("Interests").
		Find(&ret).Where("campaign_id = ?", campaignId).
		Where("start <= ?", time.Now()).Where("end >= ?", time.Now()).Error

	if err != nil{
		util.Tracer.LogError(span, err)
	}

	return ret, err
}

func (repo *CampaignRepository) GetCampaignById(ctx context.Context, campaignId uint) (model.Campaign, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCampaignById-repository")
	defer util.Tracer.FinishSpan(span)

	var ret model.Campaign

	err := repo.Database.Preload("CampaignParameters.Interests").
		Preload("CampaignParameters.CampaignRequests").
		Preload("CampaignParameters.Timestamps").
		Find(&ret).Where("id = ?", campaignId).Error

	if err != nil{
		util.Tracer.LogError(span, err)
	}

	return ret, err
}

func (repo *CampaignRepository) GetLastActiveParametersForCampaign(ctx context.Context, campaignId uint) (model.CampaignParameters, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetLastActiveParametersForCampaign-repository")
	defer util.Tracer.FinishSpan(span)

	var ret model.CampaignParameters
	result := repo.Database.Preload("Interests").Preload("CampaignRequests").Preload("Timestamps").
		Raw("select * from campaign_parameters p where p.campaign_id = ? AND p.end > NOW() AND deleted_at IS NULL AND p.start < NOW()", campaignId).First(&ret)
	if result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		return model.CampaignParameters{},result.Error
	}else if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("last active parameter for campaign id %v not found",campaignId))
		return model.CampaignParameters{}, gorm.ErrRecordNotFound
	}
	return ret, nil
}

func (repo *CampaignRepository) GetAllActiveParameters(ctx context.Context) ([]model.CampaignParameters, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllActiveParameters-repository")
	defer util.Tracer.FinishSpan(span)

	var ret []model.CampaignParameters
	result := repo.Database.Preload("Interests").Preload("CampaignRequests").Preload("Timestamps").
		Where("end > ? AND deleted_at IS NULL AND start < ? ", time.Now(), time.Now()).Find(&ret)

	if result.Error != nil {
		util.Tracer.LogError(span, result.Error)
		return make([]model.CampaignParameters, 0), result.Error
	}else if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("active parameters not found"))
		return make([]model.CampaignParameters, 0), gorm.ErrRecordNotFound
	}

	return ret, nil
}

func (repo *CampaignRepository) GetPostIDsFromCampaignIDs(ctx context.Context, campaignIDs []uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPostIDsFromCampaignIDs-repository")
	defer util.Tracer.FinishSpan(span)

	length := len(campaignIDs)
	if length == 0 {
		return make([]string, 0), nil
	}
	var ret []string
	var builder strings.Builder
	// https://stackoverflow.com/questions/39348610/select-query-with-in-clause-having-duplicate-values-in-in-clause
	builder.WriteString("select c.post_id from campaigns c join (")
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
		util.Tracer.LogError(span, result.Error)
		return make([]string, 0), result.Error
	}else if result.RowsAffected == 0 {
		return make([]string, 0), gorm.ErrRecordNotFound
	}
	return ret, nil
}

func (repo *CampaignRepository) UpdateCampaignRequest(ctx context.Context, id string, status model.RequestStatus) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateCampaignRequest-repository")
	defer util.Tracer.FinishSpan(span)

	if res := repo.Database.Model(&model.CampaignRequest{}).Where("id = ?", id).Update("request_status", status); res.Error != nil{
		util.Tracer.LogError(span, res.Error)
		return res.Error
	}else if res.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("campaign request not found"))
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *CampaignRepository) GetCampaignRequestInfluencerId(ctx context.Context, requestId uint) uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCampaignRequestInfluencerId-repository")
	defer util.Tracer.FinishSpan(span)

	var ret uint
	res := repo.Database.Raw("select influencer_id from campaign_requests " +
		"where id = ? and request_status = 0", requestId).Scan(&ret)
	if res.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("influencer not found"))
		return 0
	}
	return ret
}
