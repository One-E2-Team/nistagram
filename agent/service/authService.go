package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
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
		Roles: []model.Role{*role}, IsValidated: false, ValidationExpire: time.Now().Add(1 * time.Hour), APIToken: util.EncryptedString{}}
	err = service.AuthRepository.CreateUser(&user)
	if err != nil {
		return err
	}
	agentHost, agentPort := util.GetAgentHostAndPort()
	message := "Visit this link in the next 60 minutes to validate your account: " + util.GetAgentProtocol() +
		"://" + agentHost + ":" + agentPort + "/validate/" + util.Uint2String(user.ID) + "/" + user.ValidationUid //+ "/" + uid
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

func (service *AuthService) ValidateUser(id string, uuid string) error {
	user, err := service.AuthRepository.GetUserByProfileID(util.String2Uint(id))
	if err != nil {
		return err
	}
	if user.ValidationUid != uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
		return fmt.Errorf("BAD_VALIDATION_DATA")
	}
	user.IsValidated = true
	user.ValidationUid = ""
	user.ValidationExpire = time.Now()
	err = service.AuthRepository.UpdateUser(*user)
	return err
}

func (service *AuthService) CreateAPIToken(apiToken string, loggedUserID uint) error {
	user, err := service.AuthRepository.GetUserByProfileID(loggedUserID)
	if err != nil {
		return err
	}
	user.APIToken.Data = apiToken
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	type LoginData struct {
		Email    string `json:"email"`
		ApiToken string `json:"apiToken"`
	}
	req := LoginData{
		Email:    user.Email,
		ApiToken: apiToken,
	}
	jsonReq, _ := json.Marshal(req)
	nistagramHost, nistagramPort := util.GetNistagramHostAndPort()
	resp, err := util.NistagramRequest(http.MethodPost, util.GetNistagramProtocol() + "://" +
		nistagramHost + ":" + nistagramPort + "/agent-api/auth/login/apitoken",
		jsonReq, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	util.SetJwt(string(body))
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