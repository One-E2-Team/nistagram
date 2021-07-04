package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"nistagram/auth/model"
	"nistagram/util"
)

type AuthRepository struct {
	Database *gorm.DB
}

func (repo *AuthRepository) CreateUser(ctx context.Context,user *model.User) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateUser-repository")
	defer util.Tracer.FinishSpan(span)

	result := repo.Database.Create(user)
	if result.RowsAffected == 0 {
		util.Tracer.LogError(span, fmt.Errorf("user not created"))
		return fmt.Errorf("user not created")
	}
	fmt.Println("User Created")
	return nil
}

func (repo *AuthRepository) GetUserByEmail(ctx context.Context,email string) (*model.User, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetUserByEmail-repository")
	defer util.Tracer.FinishSpan(span)

	user := &model.User{}
	if err := repo.Database.Preload("Roles").Table("users").First(&user, "email = ?", email).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepository) GetUserByProfileID(ctx context.Context,id uint) (*model.User, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetUserByProfileID-repository")
	defer util.Tracer.FinishSpan(span)

	user := &model.User{}
	if err := repo.Database.Table("users").First(&user, "profile_id = ?", id).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepository) UpdateUser(ctx context.Context,user model.User) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateUser-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.Database.Save(user).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *AuthRepository) DeleteUser(ctx context.Context, user *model.User) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteUser-repository")
	defer util.Tracer.FinishSpan(span)

	if err := repo.Database.Delete(user).Error; err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (repo *AuthRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetRoleByName-repository")
	defer util.Tracer.FinishSpan(span)

	role := &model.Role{}
	if err := repo.Database.Table("roles").First(&role, "name = ?", name).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return role, nil
}

func (repo *AuthRepository) GetPrivilegesByUserID(ctx context.Context, id uint) (*[]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPrivilegesByUserID-repository")
	defer util.Tracer.FinishSpan(span)

	var privileges []string
	if err := repo.Database.Raw("select p.name from privileges p, role_privileges rp "+
		"where rp.role_id in (select r.id from roles r, user_roles ur where ur.user_id = ? and ur.role_id = r.id) "+
		"and p.id = rp.privilege_id", id).Scan(&privileges).Error; err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return &privileges, nil
}
