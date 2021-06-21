package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"nistagram/profile/model"
)

type ProfileRepository struct {
	RelationalDatabase *gorm.DB
	Client             *redis.Client
	Context            context.Context
}

func (repo *ProfileRepository) CreateProfile(profile *model.Profile) error {
	result := repo.RelationalDatabase.Create(profile)
	if result.RowsAffected == 0 {
		return fmt.Errorf("ProfileVertex not created")
	}
	fmt.Println("ProfileVertex Created")
	return nil
}

func (repo *ProfileRepository) FindInterestByName(name string) model.Interest {
	interest := &model.Interest{}
	repo.RelationalDatabase.Find(&interest, "name", name)
	return *interest
}

func (repo *ProfileRepository) FindUsernameContains(username string) []string {
	var result []string
	repo.RelationalDatabase.Table("profiles").Select("username").Find(&result, "username LIKE ?", "%"+username+"%")
	return result
}

func (repo *ProfileRepository) FindProfileByUsername(username string) (*model.Profile, error) {
	profile := &model.Profile{}
	if err := repo.RelationalDatabase.Preload("ProfileSettings").Preload("PersonalData").First(&profile, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (repo *ProfileRepository) GetProfileByID(id uint) (*model.Profile, error) {
	profile := &model.Profile{}
	if err := repo.RelationalDatabase.Preload("ProfileSettings").Preload("PersonalData").First(&profile, "ID = ?", id).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (repo *ProfileRepository) UpdateProfile(profile *model.Profile) error {
	if err := repo.RelationalDatabase.Save(profile).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdateProfileSettings(settings model.ProfileSettings) error {
	if err := repo.RelationalDatabase.Save(settings).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdatePersonalData(personalData model.PersonalData) error {
	if err := repo.RelationalDatabase.Save(personalData).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetAllInterests() ([]string, error) {
	var interests []string
	result := repo.RelationalDatabase.Table("interests").Select("name").Find(&interests, "name LIKE ?", "%%")
	return interests, result.Error
}

func (repo *ProfileRepository) GetCategoryByName(name string) (*model.Category, error) {
	category := &model.Category{}
	result := repo.RelationalDatabase.Table("categories").First(&category, "name = ?", name)
	return category, result.Error
}

func (repo *ProfileRepository) GetAllCategories() ([]string, error) {
	var categories []string
	result := repo.RelationalDatabase.Table("categories").Select("name").Find(&categories)
	return categories, result.Error
}

func (repo *ProfileRepository) GetVerificationRequests() ([]model.VerificationRequest, error) {
	var requests []model.VerificationRequest
	if err := repo.RelationalDatabase.Preload("Category").Table("verification_requests").Find(&requests, "verification_status = 0").Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (repo *ProfileRepository) CreateVerificationRequest(verReq *model.VerificationRequest) error {
	result := repo.RelationalDatabase.Create(verReq)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Verification request not created")
	}
	fmt.Println("Verification request created")
	return nil
}

func (repo *ProfileRepository) GetVerificationRequestById(id uint) (*model.VerificationRequest, error) {
	request := &model.VerificationRequest{}
	if err := repo.RelationalDatabase.Table("verification_requests").First(&request, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *ProfileRepository) UpdateVerificationRequest(request model.VerificationRequest) error {
	if err := repo.RelationalDatabase.Save(request).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) DeleteVerificationRequest(request *model.VerificationRequest) error {
	if err := repo.RelationalDatabase.Delete(request).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) DeleteProfile(profileId uint) error {
	profile, err := repo.GetProfileByID(profileId)
	if err != nil {
		return err
	}
	return repo.RelationalDatabase.Delete(profile).Error
}

func (repo *ProfileRepository) SendAgentRequest(request *model.AgentRequest) error {
	if err := repo.RelationalDatabase.Create(request).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetAgentRequests() ([]model.AgentRequest, error) {
	var requests []model.AgentRequest
	if err := repo.RelationalDatabase.Table("agent_requests").
		Raw("select * from agent_requests").Scan(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (repo *ProfileRepository) GetAgentRequestByProfileID(id uint) (*model.AgentRequest, error) {
	request := &model.AgentRequest{}
	if err := repo.RelationalDatabase.Table("agent_requests").First(&request, "profile_id = ?", id).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *ProfileRepository) DeleteAgentRequest(request *model.AgentRequest) error {
	if err := repo.RelationalDatabase.Unscoped().Delete(request).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetByInterests(interests []string) ([]model.Profile, error) {
	var profiles []model.Profile
	if err := repo.RelationalDatabase.Raw("select * from profiles p where p.id in " +
		"(select pd.profile_id from personal_data pd where pd.id in " +
		"(select pi.personal_data_id from person_interests pi where pi.interest_id in" +
		"(select i.id from interests i where i.name in (?))))", interests).Scan(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (repo *ProfileRepository) InsertInRedis(key string, value string) error {
	//model := model.Interest{
	//	Model: gorm.Model{},
	//	Name:  "",
	//}
	//b, err := json.Marshal(model)
	//if err != nil {
	//	return err
	//}
	_, err := repo.GetFromRedis(key)
	if err != nil {
		statusCmd := repo.Client.Set(repo.Context, key, value, 0)
		if statusCmd.Err() != nil {
			return statusCmd.Err()
		}
		return nil
	}
	return fmt.Errorf("KEY: '" + key + "' EXISTS")

}

func (repo *ProfileRepository) GetFromRedis(key string) (string, error) {
	value := repo.Client.Get(repo.Context, key)
	if value.Err() != nil {
		return "", value.Err()
	}
	return value.Val(), nil
}
