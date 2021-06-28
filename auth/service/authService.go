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

func (service *AuthService) AgentLoginAPIToken(data dto.APILoginDTO) (*dto.TokenResponseDTO, error) {
	user, err := service.AuthRepository.GetUserByEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("'" + data.Email + "' " + err.Error())
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.APIToken), []byte(data.ApiToken)); err == nil {
		token, err := util.CreateAgentToken(user.ProfileId)
		if err != nil {
			return nil, fmt.Errorf("'" + data.Email + "' " + err.Error())
		}
		resp := dto.TokenResponseDTO{
			Token:     token,
			Email:     user.Email,
			Username:  user.Username,
			ProfileId: user.ProfileId,
			Roles:     user.Roles,
		}
		return &resp, nil
	} else {
		return nil, fmt.Errorf("'" + data.Email + "' submitted invalid api token")
	}
}

func (service *AuthService) LogIn(dto dto.LogInDTO) (*model.User, error) {
	user, err := service.AuthRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, fmt.Errorf("'" + dto.Email + "' " + err.Error())
	}
	if !user.IsValidated {
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " NOT VALIDATED")
	}
	if user.IsDeleted {
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " DELETED")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " " + err.Error())
	}
	/*if !util.ValidateTOTP(user.TotpUrl.Data, dto.Passcode) {
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " entered wrong passcode.")
	}*/
	return user, nil
}

func (service *AuthService) Register(dto dto.RegisterDTO) error {
	pass := hashAndSalt(dto.Password)
	role, err := service.AuthRepository.GetRoleByName("REGULAR")
	if err != nil {
		return err
	}
	/*totpUrl, img, err := util.GenerateTOTP(dto.Email)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	png.Encode(&buf, *img)

	uid := uuid.NewString()

	file, err := os.OpenFile("../../nistagramstaticdata/totp/"+uid+".png", os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	if err != nil {
		return err
	}

	n, err := file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("Bytes written: ", n)*/

	user := model.User{ProfileId: util.String2Uint(dto.ProfileIdString), Password: pass /*TotpUrl: util.EncryptedString{Data: totpUrl}, */, Email: dto.Email, Username: dto.Username,
		ValidationUid: uuid.NewString(), Roles: []model.Role{*role}, IsDeleted: false, IsValidated: false, ValidationExpire: time.Now().Add(1 * time.Hour)}
	err = service.AuthRepository.CreateUser(&user)
	if err != nil {
		return err
	}
	//TODO: change host, port and html page
	gatewayHost, gatewayPort := util.GetGatewayHostAndPort()
	message := "Visit this link in the next 60 minutes to validate your account: " + util.GetMicroservicesProtocol() +
		"://" + gatewayHost + ":" + gatewayPort + "/api/auth/validate/" + util.Uint2String(user.ProfileId) + "/" + user.ValidationUid //+ "/" + uid
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
		return fmt.Errorf("NOT VALIDATED")
	}
	user.ValidationUid = uuid.NewString()
	user.ValidationExpire = time.Now().Add(20 * time.Minute)
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	frontHost, frontPort := util.GetFrontHostAndPort()
	var message = "Visit this link in the next 20 minutes to change your password: " + util.GetFrontProtocol() + "://" + frontHost + ":" + frontPort + "/web#/reset-password?id=" +
		util.Uint2String(user.ProfileId) + "&str=" + user.ValidationUid
	go util.SendMail(user.Email, "Recovery password", message)
	return nil
}

func (service *AuthService) ChangePassword(dto dto.RecoveryDTO) error {
	user, err := service.AuthRepository.GetUserByProfileID(util.String2Uint(dto.Id))
	if err != nil {
		return err
	}
	if !user.IsValidated {
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " NOT VALIDATED")
	}
	if user.IsDeleted {
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " DELETED")
	}
	if user.ValidationUid != dto.Uuid {
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " BAD UUID")
	}
	if time.Now().Unix() > user.ValidationExpire.Unix() {
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " VALIDATION DATE EXPIRED")
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
	user.Username = dto.Username
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

func (service *AuthService) GetAgentAPIToken(loggedUserID uint) (string, error) {
	user, err := service.AuthRepository.GetUserByProfileID(loggedUserID)
	if err != nil {
		return "", err
	}
	apiToken := uuid.NewString()
	user.APIToken = hashAndSalt(apiToken)
	role, err := service.AuthRepository.GetRoleByName("AGENT_API_CLIENT")
	if err != nil {
		return "", err
	}
	var hasRole = false
	for _, val := range user.Roles {
		if val.Name == "AGENT_API_CLIENT" {
			hasRole = true
		}
	}
	if !hasRole {
		user.Roles = append(user.Roles, *role)
	}
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return "", err
	}
	return apiToken, nil
}

func (service *AuthService) GetPrivileges(id uint) *[]string {
	user, err := service.AuthRepository.GetUserByProfileID(id)
	if err != nil {
		return nil
	}
	privileges, err := service.AuthRepository.GetPrivilegesByUserID(user.ID)
	if err != nil {
		return nil
	}
	return privileges
}

func (service *AuthService) BanUser(profileID uint) error {
	user, err := service.AuthRepository.GetUserByProfileID(profileID)
	if err != nil {
		return err
	}
	err = service.AuthRepository.DeleteUser(user)
	if err != nil {
		return err
	}
	message := "Unfortunately, your account has been deleted from Nistagram due inappropriate posts!"
	go util.SendMail(user.Email, "Deleted account", message)
	return nil
}

func (service *AuthService) MakeAgent(profileID uint) error {
	user, err := service.AuthRepository.GetUserByProfileID(profileID)
	if err != nil {
		return err
	}
	role, err := service.AuthRepository.GetRoleByName("AGENT")
	if err != nil {
		return err
	}
	user.Roles = append(user.Roles, *role)
	return service.AuthRepository.UpdateUser(*user)
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
