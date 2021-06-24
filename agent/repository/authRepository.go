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

func (repo *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := repo.Database.Preload("Roles").Table("users").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepository) GetPrivilegesByUserID(id uint) (*[]string, error) {
	var privileges []string
	if err := repo.Database.Raw("select p.name from privileges p, role_privileges rp "+
		"where rp.role_id in (select r.id from roles r, user_roles ur where ur.user_id = ? and ur.role_id = r.id) "+
		"and p.id = rp.privilege_id", id).Scan(&privileges).Error; err != nil {
		return nil, err
	}
	return &privileges, nil
}