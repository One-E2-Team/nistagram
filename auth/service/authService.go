package service

import (
	"errors"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/repository"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) LogIn(dto dto.LogInDTO) error {
	credentials := model.Credentials{}
	service.AuthRepository.Database.Find(&credentials, "username", dto.Username)
	if credentials.Password != dto.Password{
		return errors.New("BAD")
	}
	return nil
}
