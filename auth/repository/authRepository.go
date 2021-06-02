package repository

import (
	"fmt"
	"gorm.io/gorm"
	"nistagram/auth/model"
)

type AuthRepository struct {
	Database *gorm.DB
}

func (repo *AuthRepository) CreateUser(user *model.User) error {
	result := repo.Database.Create(user)
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User not created")
	}
	fmt.Println("User Created")
	return nil
}

func (repo *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := repo.Database.First(&user, "username = ?", username).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepository) GetUserByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := repo.Database.First(&user, "ID = ?", id).Error; err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepository) UpdateUser(user *model.User) {
	repo.Database.Save(user)
}
