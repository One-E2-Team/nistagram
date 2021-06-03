package service

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/repository"
	"nistagram/util"
	"time"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) LogIn(dto dto.LogInDTO) (*model.User, error) {
	user, err := service.AuthRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted {
		return nil, fmt.Errorf("USER_DELETED")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) Register(dto dto.RegisterDTO) error {
	pass := hashAndSalt(dto.Password)
	user := model.User{ProfileId: util.String2Uint(dto.ProfileIdString), Password: pass,
		Email: dto.Email, Roles: nil, IsDeleted: false, ValidationExpire: time.Now()}
	err := service.AuthRepository.CreateUser(&user)
	return err
}

func (service *AuthService) GetUserByEmail(email string) *model.User {
	user, err := service.AuthRepository.GetUserByEmail(email)
	if err != nil {
		return nil
	}
	return user
}

func (service *AuthService) RequestPassRecovery(email string) error {
	var user *model.User
	user, err := service.AuthRepository.GetUserByEmail(email)
	if err != nil {
		return err
	}
	user.ValidationUid = uuid.NewString()
	user.ValidationExpire = time.Now().Add(20 * time.Minute)
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	//TODO: change host, port and html page
	var message = "Visit this link in the next 20 minutes to change your password: http://localhost:8000/recovery.html?id=" +
		util.Uint2String(user.ProfileId) + "&str=" + user.ValidationUid
	go util.SendMail(user.Email, "Recovery password", message)
	return nil
}

func (service *AuthService) ChangePassword(dto dto.RecoveryDTO) error {
	user, err := service.AuthRepository.GetUserByProfileID(dto.Id)
	if err != nil {
		return err
	}
	if user.IsDeleted || user.ValidationUid != dto.Uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
		return fmt.Errorf("BAD_RECOVERY_DATA")
	}
	user.ValidationExpire = time.Now()
	user.ValidationUid = ""
	user.Password = hashAndSalt(dto.Password)
	err = service.AuthRepository.UpdateUser(*user)
	return err
}

func (service *AuthService) UpdateUser(dto dto.UpdateUserDTO) error {
	user, err := service.AuthRepository.GetUserByProfileID(util.String2Uint(dto.ProfileId))
	if err != nil {
		return err
	}
	user.Email = dto.Email
	err = service.AuthRepository.UpdateUser(*user)
	return err
}

func hashAndSalt(pass string) string {
	bytePass := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(hash)
}
