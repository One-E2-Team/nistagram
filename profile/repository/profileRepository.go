package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"nistagram/profile/model"
	"nistagram/util"
)

type ProfileRepository struct {
	RelationalDatabase *gorm.DB
	Client             *redis.Client
	Context            context.Context
}

func (repo *ProfileRepository) CreateProfile(ctx context.Context, profile *model.Profile) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateProfile-repository")
	defer util.Tracer.FinishSpan(span)

	result := repo.RelationalDatabase.Create(profile)
	if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("ProfileVertex not created"))
		return fmt.Errorf("ProfileVertex not created")
	}
	fmt.Println("ProfileVertex Created")
	return nil
}

func (repo *ProfileRepository) FindInterestByName(ctx context.Context, name string) model.Interest {
	span := util.Tracer.StartSpanFromContext(ctx, "FindInterestByName-repository")
	defer util.Tracer.FinishSpan(span)

	interest := &model.Interest{}
	repo.RelationalDatabase.Find(&interest, "name", name)
	return *interest
}

func (repo *ProfileRepository) FindUsernameContains(ctx context.Context, username string) []string {
	span := util.Tracer.StartSpanFromContext(ctx, "FindUsernameContains-repository")
	defer util.Tracer.FinishSpan(span)

	var result []string
	repo.RelationalDatabase.Table("profiles").Select("username").Find(&result, "username LIKE ?", "%"+username+"%")
	return result
}

func (repo *ProfileRepository) FindProfileByUsername(ctx context.Context, username string) (*model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "FindProfileByUsername-repository")
	defer util.Tracer.FinishSpan(span)

	profile := &model.Profile{}
	if err := repo.RelationalDatabase.Preload("ProfileSettings").Preload("PersonalData").First(&profile, "username = ?", username).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (repo *ProfileRepository) GetProfileByID(ctx context.Context, id uint) (*model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileByID-repository")
	defer util.Tracer.FinishSpan(span)
	profile := &model.Profile{}
	if err := repo.RelationalDatabase.Preload("ProfileSettings").Preload("PersonalData").First(&profile, "ID = ?", id).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (repo *ProfileRepository) UpdateProfile(ctx context.Context, profile *model.Profile) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateProfile-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Save(profile).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdateProfileSettings(ctx context.Context, settings model.ProfileSettings) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateProfileSettings-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Save(settings).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdatePersonalData(ctx context.Context, personalData model.PersonalData) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdatePersonalData-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Save(personalData).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetAllInterests(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllInterests-repository")
	defer util.Tracer.FinishSpan(span)

	var interests []string
	result := repo.RelationalDatabase.Table("interests").Select("name").Find(&interests, "name LIKE ?", "%%")
	return interests, result.Error
}

func (repo *ProfileRepository) GetCategoryByName(ctx context.Context, name string) (*model.Category, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCategoryByName-repository")
	defer util.Tracer.FinishSpan(span)

	category := &model.Category{}
	result := repo.RelationalDatabase.Table("categories").First(&category, "name = ?", name)
	return category, result.Error
}

func (repo *ProfileRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllCategories-repository")
	defer util.Tracer.FinishSpan(span)

	var categories []string
	result := repo.RelationalDatabase.Table("categories").Select("name").Find(&categories)
	return categories, result.Error
}

func (repo *ProfileRepository) GetVerificationRequests(ctx context.Context) ([]model.VerificationRequest, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetVerificationRequests-repository")
	defer util.Tracer.FinishSpan(span)

	var requests []model.VerificationRequest
	if err := repo.RelationalDatabase.Preload("Category").Table("verification_requests").Find(&requests, "verification_status = 0").Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return requests, nil
}

func (repo *ProfileRepository) CreateVerificationRequest(ctx context.Context, verReq *model.VerificationRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateVerificationRequest-repository")
	defer util.Tracer.FinishSpan(span)

	result := repo.RelationalDatabase.Create(verReq)
	if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("verification request not created"))
		return fmt.Errorf("verification request not created")
	}
	fmt.Println("Verification request created")
	return nil
}

func (repo *ProfileRepository) GetVerificationRequestById(ctx context.Context, id uint) (*model.VerificationRequest, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetVerificationRequestById-repository")
	defer util.Tracer.FinishSpan(span)

	request := &model.VerificationRequest{}
	if err := repo.RelationalDatabase.Table("verification_requests").First(&request, "id = ?", id).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return request, nil
}

func (repo *ProfileRepository) UpdateVerificationRequest(ctx context.Context, request model.VerificationRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateVerificationRequest-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Save(request).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) DeleteVerificationRequest(ctx context.Context, request *model.VerificationRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteVerificationRequest-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Delete(request).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) DeleteProfile(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteProfile-repository")
	defer util.Tracer.FinishSpan(span)

	profile, err := repo.GetProfileByID(context.Background(), profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return repo.RelationalDatabase.Delete(profile).Error
}

func (repo *ProfileRepository) SendAgentRequest(ctx context.Context, request *model.AgentRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "SendAgentRequest-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Create(request).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetAgentRequests(ctx context.Context) ([]model.AgentRequest, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAgentRequests-repository")
	defer util.Tracer.FinishSpan(span)

	var requests []model.AgentRequest
	if err := repo.RelationalDatabase.Table("agent_requests").
		Raw("select * from agent_requests").Scan(&requests).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return requests, nil
}

func (repo *ProfileRepository) GetAgentRequestByProfileID(ctx context.Context, id uint) (*model.AgentRequest, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAgentRequestByProfileID-repository")
	defer util.Tracer.FinishSpan(span)

	request := &model.AgentRequest{}
	if err := repo.RelationalDatabase.Table("agent_requests").First(&request, "profile_id = ?", id).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return request, nil
}

func (repo *ProfileRepository) DeleteAgentRequest(ctx context.Context, request *model.AgentRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteAgentRequest-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.RelationalDatabase.Unscoped().Delete(request).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetByInterests(ctx context.Context, interests []string) ([]model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetByInterests-repository")
	defer util.Tracer.FinishSpan(span)

	var profiles []model.Profile
	if err := repo.RelationalDatabase.Raw("select * from profiles p where p.id in "+
		"(select pd.profile_id from personal_data pd where pd.id in "+
		"(select pi.personal_data_id from person_interests pi where pi.interest_id in"+
		"(select i.id from interests i where i.name in (?))))", interests).Scan(&profiles).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profiles, nil
}

func (repo *ProfileRepository) GetProfileIdsByUsernames(ctx context.Context, usernames []string) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileIdsByUsernames-repository")
	defer util.Tracer.FinishSpan(span)

	var ret []string
	if err := repo.RelationalDatabase.Table("profiles").Raw("select p.id from profiles p where p.username in (?)", usernames).Scan(&ret).Error; err != nil {
		util.Tracer.LogError(span, err)
		return make([]string, 0), err
	}
	return ret, nil
}

func (repo *ProfileRepository) GetProfileInterests(ctx context.Context, id uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileInterests-repository")
	defer util.Tracer.FinishSpan(span)

	var ret []string
	if err := repo.RelationalDatabase.Raw("select i.name from interests i where i.id in" +
		"(select pi.interest_id from person_interests pi where pi.personal_data_id = ?)", id).Scan(&ret).Error; err != nil {
		util.Tracer.LogError(span, err)
		return make([]string, 0), err
	}
	return ret, nil
}

func (repo *ProfileRepository) InsertInRedis(ctx context.Context, key string, value string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "InsertInRedis-repository")
	defer util.Tracer.FinishSpan(span)
	//model := model.Interest{
	//	Model: gorm.Model{},
	//	Name:  "",
	//}
	//b, err := json.Marshal(model)
	//if err != nil {
	//	return err
	//}
	_, err := repo.GetFromRedis(ctx, key)
	if err != nil {
		statusCmd := repo.Client.Set(repo.Context, key, value, 0)
		if statusCmd.Err() != nil {
			util.Tracer.LogError(span, err)
			return statusCmd.Err()
		}
		return nil
	}
	return fmt.Errorf("KEY: '" + key + "' EXISTS")

}

func (repo *ProfileRepository) GetFromRedis(ctx context.Context, key string) (string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetFromRedis-repository")
	defer util.Tracer.FinishSpan(span)

	value := repo.Client.Get(repo.Context, key)
	if value.Err() != nil {
		util.Tracer.LogError(span, value.Err())
		return "", value.Err()
	}
	return value.Val(), nil
}

func (repo *ProfileRepository) GetPersonalDataByProfileId(ctx context.Context, id uint) (*model.PersonalData, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPersonalDataByProfileId-repository")
	defer util.Tracer.FinishSpan(span)

	data := &model.PersonalData{}
	if err := repo.RelationalDatabase.Table("personal_data").Preload("InterestedIn").First(&data, "profile_id = ?", id).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return data, nil
}
