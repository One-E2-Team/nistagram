package repository

import (
	"gorm.io/gorm"
)

type AuthRepository struct {
	Database *gorm.DB
}
