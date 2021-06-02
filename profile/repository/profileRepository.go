package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/profile/model"
)

type ProfileRepository struct {
	Database *gorm.DB
}

func (repo *ProfileRepository) CreateProfile(profile *model.Profile) error{
	result := repo.Database.Create(profile)
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Profile not created")
	}
	fmt.Println("Profile Created")
	return nil
}

func (repo *ProfileRepository) FindInterestByName(name string) model.Interest{
	interest := &model.Interest{}
	repo.Database.Find(&interest, "name", name)
	return *interest
}

func (repo *ProfileRepository) FindUsernameContains(username string) []string{
	var result []string
	repo.Database.Table("profiles").Select("username").Find(&result, "username LIKE ?", "%" + username + "%")
	return result
}

func (repo *ProfileRepository) FindProfileByUsername(username string) *model.Profile{
	profile := &model.Profile{}
	repo.Database.Preload("PersonalData").First(&profile, "username = ?", username)
	return profile
}
