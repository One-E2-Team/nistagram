package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/profile/model"
)

type ProfileRepository struct {
	Database *gorm.DB
}

func (repo *ProfileRepository) CreateProfile(profile *model.Profile) error {
	result := repo.Database.Create(profile)
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Profile not created")
	}
	fmt.Println("Profile Created")
	return nil
}

func (repo *ProfileRepository) FindInterestByName(name string) model.Interest {
	interest := &model.Interest{}
	repo.Database.Find(&interest, "name", name)
	return *interest
}

func (repo *ProfileRepository) FindUsernameContains(username string) []string {
	var result []string
	repo.Database.Table("profiles").Select("username").Find(&result, "username LIKE ?", "%"+username+"%")
	return result
}

func (repo *ProfileRepository) FindProfileByUsername(username string) *model.Profile {
	profile := &model.Profile{}
	repo.Database.Preload("PersonalData").First(&profile, "username = ?", username)
	return profile
}

func (repo *ProfileRepository) GetProfileByID(id uint) (*model.Profile, error) {
	profile := &model.Profile{}
	if err := repo.Database.Preload("ProfileSettings").Preload("PersonalData").First(&profile, "ID = ?", id).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (repo *ProfileRepository) UpdateProfile(profile *model.Profile) error {
	if err := repo.Database.Save(profile).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdateProfileSettings(settings model.ProfileSettings) error {
	if err := repo.Database.Save(settings).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProfileRepository) UpdatePersonalData(personalData model.PersonalData) error {
	if err := repo.Database.Save(personalData).Error; err != nil {
		return err
	}
	return nil
}
