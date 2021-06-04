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

func (repo *ProfileRepository) FindProfileByUsername(username string) *model.Profile {
	profile := &model.Profile{}
	repo.RelationalDatabase.Preload("PersonalData").First(&profile, "username = ?", username)
	return profile
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
