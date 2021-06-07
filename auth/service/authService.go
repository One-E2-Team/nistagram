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
	if !user.IsValidated {
		return nil, fmt.Errorf("USER_NOT_VALIDATED")
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
	user := model.User{ProfileId: util.String2Uint(dto.ProfileIdString), Password: pass, Email: dto.Email,
		ValidationUid: uuid.NewString(), Roles: nil, IsDeleted: false, IsValidated: false, ValidationExpire: time.Now().Add(1 * time.Hour)}
	err := service.AuthRepository.CreateUser(&user)
	if err != nil {
		return err
	}
	//TODO: change host, port and html page
	message := "Visit this link in the next 60 minutes to validate your account: http://localhost:81/api/auth/validate/" +
	util.Uint2String(user.ProfileId) + "/" + user.ValidationUid
	go util.SendMail(dto.Email, "Account Validation", message)
	return nil
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
	if !user.IsValidated {
		return fmt.Errorf("USER_NOT_VALIDATED")
	}
	user.ValidationUid = uuid.NewString()
	user.ValidationExpire = time.Now().Add(20 * time.Minute)
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	//TODO: change host, port
	var message = "Visit this link in the next 20 minutes to change your password: http://localhost:3000/#/reset-password?id=" +
		util.Uint2String(user.ProfileId) + "&str=" + user.ValidationUid
	go util.SendMail(user.Email, "Recovery password", message)
	return nil
}

func (service *AuthService) ChangePassword(dto dto.RecoveryDTO) error {
	user, err := service.AuthRepository.GetUserByProfileID(util.String2Uint(dto.Id))
	if err != nil {
		return err
	}
	if !user.IsValidated || user.IsDeleted || user.ValidationUid != dto.Uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
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

func (service *AuthService) ValidateUser(id string, uuid string) error {
	user, err := service.AuthRepository.GetUserByProfileID(util.String2Uint(id))
	if err != nil {
		return err
	}
	if user.IsDeleted || user.ValidationUid != uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
		return fmt.Errorf("BAD_VALIDATION_DATA")
	}
	user.IsValidated = true
	user.ValidationUid = ""
	user.ValidationExpire = time.Now()
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
