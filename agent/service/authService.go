package service

import(
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"nistagram/util"
	"time"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) Register(dto dto.RegisterDTO) error {
	pass := hashAndSalt(dto.Password)
	role, err := service.AuthRepository.GetRoleByName("CUSTOMER")
	if err != nil {
		return err
	}

	user := model.User{Email: dto.Email, Password: pass,Address: dto.Address, ValidationUid: uuid.NewString(),
		Roles: []model.Role{*role}, IsValidated: false, ValidationExpire: time.Now().Add(1 * time.Hour)}
	err = service.AuthRepository.CreateUser(&user)
	if err != nil {
		return err
	}
	gatewayHost, gatewayPort := util.GetGatewayHostAndPort()
	message := "Visit this link in the next 60 minutes to validate your account: " + util.GetMicroservicesProtocol() +
		"://" + gatewayHost + ":" + gatewayPort + "/api/auth/validate/" + util.Uint2String(user.ID) + "/" + user.ValidationUid //+ "/" + uid
	go util.SendMail(dto.Email, "Account Validation", message)
	return nil
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