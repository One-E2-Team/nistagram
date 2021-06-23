package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/agent/model"
)

type AuthRepository struct {
	Database *gorm.DB
}

func (repo *AuthRepository) CreateUser(user *model.User) error {
	result := repo.Database.Create(user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User not created")
	}
	fmt.Println("User Created")
	return nil
}

func (repo *AuthRepository) GetRoleByName(name string) (*model.Role, error) {
	role := &model.Role{}
	if err := repo.Database.Table("roles").First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return role, nil
}