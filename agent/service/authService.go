package service

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"nistagram/agent/util"
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
	agentHost, agentPort := util.GetAgentHostAndPort()
	message := "Visit this link in the next 60 minutes to validate your account: " + util.GetAgentProtocol() +
		"://" + agentHost + ":" + agentPort + "/api/auth/validate/" + util.Uint2String(user.ID) + "/" + user.ValidationUid //+ "/" + uid
	go util.SendMail(dto.Email, "Account Validation", message)
	return nil
}

func (service *AuthService) LogIn(dto dto.LogInDTO) (*model.User, error) {
	user, err := service.AuthRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("'" + dto.Email + "' " + err.Error())
	}
	if !user.IsValidated {
		return nil, fmt.Errorf(util.Uint2String(user.ID) + " NOT VALIDATED")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) GetPrivileges(id uint) *[]string {
	privileges, err := service.AuthRepository.GetPrivilegesByUserID(id)
	if err != nil {
		return nil
	}
	return privileges
}

func (service *AuthService) RBAC(handler func(http.ResponseWriter, *http.Request),privilege string, returnCollection bool) func(http.ResponseWriter, *http.Request) {

	finalHandler := func(pass bool) func(http.ResponseWriter, *http.Request) {
		if pass {
			return handler
		} else {
			return func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Header().Set("Content-Type", "application/json")
				if returnCollection {
					_, _ = writer.Write([]byte("[{\"status\":\"fail\", \"reason\":\"unauthorized\"}]"))
				} else {
					_, _ = writer.Write([]byte("{\"status\":\"fail\", \"reason\":\"unauthorized\"}"))
				}
			}
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		var handleFunc func(http.ResponseWriter, *http.Request)
		id := util.GetLoggedUserIDFromToken(request)
		if id == 0 {
			handleFunc = finalHandler(false)
		} else {
			validPrivileges := service.GetPrivileges(id)
			valid := false
			for _, val := range *validPrivileges {
				if val == privilege {
					valid = true
					break
				}
			}
			if valid {
				handleFunc = finalHandler(true)
			} else {
				handleFunc = finalHandler(false)
			}
		}
		handleFunc(writer, request)
	}
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