package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/auth/model"
)

type AuthRepository struct {
	Database *gorm.DB
}

func (repo *AuthRepository) CreateUser(user *model.User) error{
	result := repo.Database.Create(user)
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User not created")
	}
	fmt.Println("User Created")
	return nil
}
