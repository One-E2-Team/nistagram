package service

import (
	"context"
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

func (service *AuthService) AgentLoginAPIToken(ctx context.Context, data dto.APILoginDTO) (*string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "AgentLoginAPIToken-service")
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("sevicing get user by email %s\n", data.Email))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	user, err := service.AuthRepository.GetUserByEmail(nextCtx, data.Email)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, fmt.Errorf("'" + data.Email + "' " + err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.APIToken), []byte(data.ApiToken)); err == nil {
		token, err := util.CreateAgentToken(user.ProfileId)	//TODO: tracking add context to util method and send nexCnt
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, fmt.Errorf("'" + data.Email + "' " + err.Error())
		}
		return &token, nil
	} else {
		return nil, fmt.Errorf("'" + data.Email + "' submitted invalid api token")
	}
}

func (service *AuthService) LogIn(ctx context.Context, dto dto.LogInDTO) (*model.User, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "LogIn-service")
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("sevicing log in \n"))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	user, err := service.AuthRepository.GetUserByEmail(nextCtx, dto.Email)

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, fmt.Errorf("'" + dto.Email + "' " + err.Error())
	}
	if !user.IsValidated {
		util.Tracer.LogError(span, fmt.Errorf("user is not validated"))
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " NOT VALIDATED")
	}
	if user.IsDeleted {
		util.Tracer.LogError(span, fmt.Errorf("user is deleted"))
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " DELETED")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " " + err.Error())
	}
	/*if !util.ValidateTOTP(user.TotpUrl.Data, dto.Passcode) {
		return nil, fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " entered wrong passcode.")
	}*/
	return user, nil
}

func (service *AuthService) Register(ctx context.Context, dto dto.RegisterDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Register-service")
	defer util.Tracer.FinishSpan(span)

	pass := hashAndSalt(dto.Password)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	role, err := service.AuthRepository.GetRoleByName(nextCtx, "REGULAR")
	if err != nil {
		util.Tracer.LogError(span, err)
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
	err = service.AuthRepository.CreateUser(nextCtx, &user)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	/*//TODO: change host, port and html page
	gatewayHost, gatewayPort := util.GetGatewayHostAndPort()
	message := "Visit this link in the next 60 minutes to validate your account: " + util.GetMicroservicesProtocol() +
		"://" + gatewayHost + ":" + gatewayPort + "/api/auth/validate/" + util.Uint2String(user.ProfileId) + "/" + user.ValidationUid //+ "/" + uid
	go util.SendMail(dto.Email, "Account Validation", message)
*/
	util.Tracer.LogFields(span, "service", fmt.Sprintf("send verification mail to %s", dto.Email))
	return nil
}

func (service *AuthService) GetUserByEmail(ctx context.Context, email string) *model.User {
	span := util.Tracer.StartSpanFromContext(ctx, "GetUserByEmail-service")
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("sevicing get user by email \n"))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByEmail(nextCtx, email)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil
	}
	return user
}

func (service *AuthService) RequestPassRecovery(ctx context.Context, email string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "RequestPassRecovery-service")
	defer util.Tracer.FinishSpan(span)

	var user *model.User

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByEmail(nextCtx, email)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if !user.IsValidated {
		util.Tracer.LogError(span, fmt.Errorf("user is not validated"))
		return fmt.Errorf("NOT VALIDATED")
	}
	user.ValidationUid = uuid.NewString()
	user.ValidationExpire = time.Now().Add(20 * time.Minute)
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	frontHost, frontPort := util.GetFrontHostAndPort()
	var message = "Visit this link in the next 20 minutes to change your password: " + util.GetFrontProtocol() + "://" + frontHost + ":" + frontPort + "/web#/reset-password?id=" +
		util.Uint2String(user.ProfileId) + "&str=" + user.ValidationUid
	go util.SendMail(user.Email, "Recovery password", message)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("send recovery password mail to %s", user.Email))
	return nil
}

func (service *AuthService) ChangePassword(ctx context.Context, dto dto.RecoveryDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangePassword-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)


	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, util.String2Uint(dto.Id))
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if !user.IsValidated {
		util.Tracer.LogError(span, fmt.Errorf("user with id %v is not validated",user.ProfileId))
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " NOT VALIDATED")
	}
	if user.IsDeleted {
		util.Tracer.LogError(span, fmt.Errorf("user with id %v is deleted",user.ProfileId))
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " DELETED")
	}
	if user.ValidationUid != dto.Uuid {
		util.Tracer.LogError(span, fmt.Errorf("user with id %v has bad uuid",user.ProfileId))
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " BAD UUID")
	}
	if time.Now().Unix() > user.ValidationExpire.Unix() {
		util.Tracer.LogError(span, fmt.Errorf("user with id %v validation date expired",user.ProfileId))
		return fmt.Errorf(util.GetLoggingStringFromID(user.ProfileId) + " VALIDATION DATE EXPIRED")
	}
	user.ValidationExpire = time.Now()
	user.ValidationUid = ""
	user.Password = hashAndSalt(dto.Password)
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	return err
}

func (service *AuthService) UpdateUser(ctx context.Context, dto dto.UpdateUserDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateUser-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, util.String2Uint(dto.ProfileId))
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	user.Email = dto.Email
	user.Username = dto.Username
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	return err
}

func (service *AuthService) ValidateUser(ctx context.Context, id string, uuid string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ValidateUser-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, util.String2Uint(id))
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if user.IsDeleted || user.ValidationUid != uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
		util.Tracer.LogError(span, fmt.Errorf("bad validation data"))
		return fmt.Errorf("BAD_VALIDATION_DATA")
	}
	user.IsValidated = true
	user.ValidationUid = ""
	user.ValidationExpire = time.Now()
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	return err
}

func (service *AuthService) GetAgentAPIToken(ctx context.Context, loggedUserID uint) (string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAgentAPIToken-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", err
	}
	apiToken := uuid.NewString()
	user.APIToken = hashAndSalt(apiToken)
	role, err := service.AuthRepository.GetRoleByName(nextCtx, "AGENT_API_CLIENT")
	if err != nil {
		util.Tracer.LogError(span, err)
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
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", err
	}
	return apiToken, nil
}

func (service *AuthService) GetPrivileges(ctx context.Context, id uint) *[]string {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPrivileges-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil
	}
	privileges, err := service.AuthRepository.GetPrivilegesByUserID(nextCtx, user.ID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil
	}
	return privileges
}

func (service *AuthService) BanUser(ctx context.Context, profileID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "BanUser-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, profileID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	err = service.AuthRepository.DeleteUser(nextCtx, user)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	message := "Unfortunately, your account has been deleted from Nistagram due inappropriate posts!"
	go util.SendMail(user.Email, "Deleted account", message)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("send deleted account mail to %s", user.Email))
	return nil
}

func (service *AuthService) MakeAgent(ctx context.Context, profileID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "MakeAgent-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	user, err := service.AuthRepository.GetUserByProfileID(nextCtx, profileID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	role, err := service.AuthRepository.GetRoleByName(nextCtx, "AGENT")
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	user.Roles = append(user.Roles, *role)
	err = service.AuthRepository.UpdateUser(nextCtx, *user)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
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
