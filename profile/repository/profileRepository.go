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
		return fmt.Errorf("Profile not created")
	}
	fmt.Println("Profile Created")
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
	if err := repo.RelationalDatabase.Preload("PersonalData").First(&profile, "username = ?", username).Error; err != nil {
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

func (repo *ProfileRepository) GetVerificationRequests() ([]model.VerificationRequest, error){
	var requests []model.VerificationRequest
	if err := repo.RelationalDatabase.Preload("Category").Table("verification_requests").Find(&requests, "verification_status = 0").Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (repo *ProfileRepository) CreateVerificationRequest(verReq *model.VerificationRequest) error{
	result := repo.RelationalDatabase.Create(verReq)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Verification request not created")
	}
	fmt.Println("Verification request created")
	return nil
}

func (repo *ProfileRepository) GetVerificationRequestById(id uint) (*model.VerificationRequest, error){
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

func (repo *ProfileRepository) DeleteVerificationRequest(request model.VerificationRequest) error {
	if err := repo.RelationalDatabase.Delete(request).Error; err != nil {
		return err
	}
	return nil
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
