package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"nistagram/profile/dto"
	"nistagram/profile/model"
	"nistagram/profile/repository"
	"nistagram/util"
)

type ProfileService struct {
	ProfileRepository *repository.ProfileRepository
}

func (service *ProfileService) Register(ctx context.Context, dto dto.RegistrationDto) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Register-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profileSettings := model.ProfileSettings{IsPrivate: dto.IsPrivate, CanReceiveMessageFromUnknown: true, CanBeTagged: true}
	personalData := model.PersonalData{Name: dto.Name, Surname: dto.Surname, Telephone: dto.Telephone,
		Gender: dto.Gender, BirthDate: dto.BirthDate}
	for _, item := range dto.InterestedIn {
		interest := service.ProfileRepository.FindInterestByName(nextCtx, item)
		personalData.AddItem(interest)
	}
	profile := model.Profile{Username: dto.Username, Email: dto.Email, ProfileSettings: profileSettings,
		PersonalData: personalData, Biography: dto.Biography, Website: dto.WebSite, IsVerified: false}
	err := service.ProfileRepository.CreateProfile(nextCtx, &profile)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	postBody, _ := json.Marshal(map[string]string{
		"profileId": util.Uint2String(profile.ID),
		"password":  dto.Password,
		"email":     profile.Email,
		"username":  profile.Username,
	})
	//responseBody := bytes.NewBuffer(postBody)
	err = registerInAuth(nextCtx, postBody)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	go func() {
		err := registerInConnection(nextCtx, profile.ID, postBody)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println("conn bug")
			fmt.Println(err)
		}
	}()
	return nil
}

func (service *ProfileService) RegisterAgent(ctx context.Context, dto dto.RegistrationDto) error {
	span := util.Tracer.StartSpanFromContext(ctx, "RegisterAgent-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	err := service.Register(nextCtx, dto)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	profile, err := service.GetProfileByUsername(nextCtx, dto.Username)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	fmt.Println("Profile id: ", profile.ID)
	err = makeAgent(nextCtx, profile.ID)

	return err
}

func registerInAuth(ctx context.Context, postBody []byte) error {
	span := util.Tracer.StartSpanFromContext(ctx, "registerInAuth-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	authHost, authPort := util.GetAuthHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/register", postBody,
		map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return err
	}
	return nil
}

func registerInConnection(ctx context.Context, profileId uint, postBody []byte) error {
	span := util.Tracer.StartSpanFromContext(ctx, "registerInConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/profile/"+util.Uint2String(profileId), postBody,
		map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return err
	}
	return nil
}

func (service *ProfileService) Search(ctx context.Context, username string) []string {
	span := util.Tracer.StartSpanFromContext(ctx, "Search-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing username %s\n", username))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.FindUsernameContains(nextCtx, username)
}

func (service *ProfileService) SearchForTag(ctx context.Context, loggedUserId uint, username string) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "SearchForTag-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v, username %s\n", loggedUserId, username))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var ret []string
	usernames := service.ProfileRepository.FindUsernameContains(nextCtx, username)

	followingProfiles, err := getFollowingProfiles(nextCtx, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}

	for i := 0; i < len(usernames); i++ {
		profile, err := service.GetProfileByUsername(nextCtx, usernames[i])
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println("Can't get profile by username!")
		}
		if util.IsFollowed(followingProfiles, profile.ID) && profile.ProfileSettings.CanBeTagged {
			ret = append(ret, profile.Username)
		}
	}

	return ret, nil
}

func getFollowingProfiles(ctx context.Context, loggedUserId uint) ([]util.FollowingProfileDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getFollowingProfiles-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n", loggedUserId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var followingProfiles []util.FollowingProfileDTO

	err = json.NewDecoder(resp.Body).Decode(&followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	return followingProfiles, err
}

func (service *ProfileService) GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileByUsername-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing username %s\n", username))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.FindProfileByUsername(nextCtx, username)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) ChangeProfileSettings(ctx context.Context, dto dto.ProfileSettingsDTO, loggedUserId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangeProfileSettings-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id  %v\n", loggedUserId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	profileSettings := profile.ProfileSettings
	if profileSettings.IsPrivate != dto.IsPrivate {
		err = service.changePrivacyInPostService(nextCtx, dto.IsPrivate, loggedUserId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}
	profileSettings.IsPrivate = dto.IsPrivate
	profileSettings.CanBeTagged = dto.CanBeTagged
	profileSettings.CanReceiveMessageFromUnknown = dto.CanReceiveMessageFromUnknown
	err = service.ProfileRepository.UpdateProfileSettings(nextCtx, profileSettings)
	return err
}

func (service *ProfileService) ChangePersonalData(ctx context.Context, dto dto.PersonalDataDTO, loggedUserId uint) (string, string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangePersonalData-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n",loggedUserId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", "", err
	}
	oldUsername, oldEmail := "", ""
	callAuth := false
	callPost := false
	if profile.Email != dto.Email {
		oldEmail = profile.Email
		callAuth = true
	}
	if profile.Username != dto.Username {
		//TODO: change data in other ms
		oldUsername = profile.Username
		callAuth = true
		callPost = true
	}
	profile.Username = dto.Username
	profile.Website = dto.Website
	profile.Biography = dto.Biography
	profile.Email = dto.Email
	profile.PersonalData.Name = dto.Name
	profile.PersonalData.BirthDate = dto.BirthDate
	profile.PersonalData.Gender = dto.Gender
	profile.PersonalData.Surname = dto.Surname
	profile.PersonalData.Telephone = dto.Telephone
	err = service.ProfileRepository.UpdateProfile(nextCtx, profile)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", "", err
	}
	if callAuth {
		postBody, _ := json.Marshal(map[string]string{
			"profileId": util.Uint2String(profile.ID),
			"email":     profile.Email,
			"username":  profile.Username,
		})
		authHost, authPort := util.GetAuthHostAndPort()
		_, err = util.CrossServiceRequest(nextCtx, http.MethodPost,
			util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/update-user", postBody,
			map[string]string{"Content-Type": "application/json;"})
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return "", "", err
		}
	}
	if callPost {
		err = service.changeUsernameInPostService(nextCtx, loggedUserId, dto.Username)
		if err != nil {
			util.Tracer.LogError(span, err)
			return "", "", err
		}
	}
	err = service.ProfileRepository.UpdatePersonalData(nextCtx, profile.PersonalData)
	return oldUsername, oldEmail, err
}

func (service *ProfileService) GetAllInterests(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllInterests-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	interests, err := service.ProfileRepository.GetAllInterests(nextCtx)
	return interests, err
}

func (service *ProfileService) GetAllCategories(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllCategories-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	categories, err := service.ProfileRepository.GetAllCategories(nextCtx)
	return categories, err
}

func (service *ProfileService) CreateVerificationRequest(ctx context.Context, profileId uint, requestDTO dto.VerificationRequestDTO, fileName string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateVerificationRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	category, err := service.ProfileRepository.GetCategoryByName(nextCtx, requestDTO.Category)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return err
	}
	var verReq = model.VerificationRequest{ProfileID: profileId, Name: requestDTO.Name, Surname: requestDTO.Surname,
		VerificationStatus: model.SENT, ImagePath: fileName, Category: *category}
	err = service.ProfileRepository.CreateVerificationRequest(nextCtx, &verReq)
	return err
}

func (service *ProfileService) UpdateVerificationRequest(ctx context.Context, verifyDTO dto.VerifyDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateVerificationRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing verification request with id %v\n", verifyDTO.VerificationId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	request, err := service.ProfileRepository.GetVerificationRequestById(nextCtx, verifyDTO.VerificationId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if verifyDTO.Status {
		request.VerificationStatus = model.VERIFIED
		err = service.ProfileRepository.UpdateVerificationRequest(nextCtx, *request)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		profile, err := service.ProfileRepository.GetProfileByID(nextCtx, request.ProfileID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		profile.IsVerified = true
		err = service.ProfileRepository.UpdateProfile(nextCtx, profile)
	} else {
		err = service.ProfileRepository.DeleteVerificationRequest(nextCtx, request)
	}

	return err
}

func (service *ProfileService) GetMyProfileSettings(ctx context.Context, loggedUserId uint) (dto.ProfileSettingsDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyProfileSettings-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n", loggedUserId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	ret := dto.ProfileSettingsDTO{}
	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}
	ret.CanReceiveMessageFromUnknown = profile.ProfileSettings.CanReceiveMessageFromUnknown
	ret.CanBeTagged = profile.ProfileSettings.CanBeTagged
	ret.IsPrivate = profile.ProfileSettings.IsPrivate
	return ret, nil
}

func (service *ProfileService) GetMyPersonalData(ctx context.Context, loggedUserId uint) (dto.PersonalDataDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyPersonalData-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n", loggedUserId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return dto.PersonalDataDTO{}, err
	}
	personalData := profile.PersonalData
	ret := dto.PersonalDataDTO{Username: profile.Username, Name: personalData.Name, Surname: personalData.Surname,
		Email: profile.Email, Telephone: personalData.Telephone, Gender: personalData.Gender,
		BirthDate: personalData.BirthDate, Biography: personalData.BirthDate, Website: profile.Website}
	return ret, nil
}

func (service *ProfileService) GetProfileByID(ctx context.Context, id uint) (*model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileByID-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) GetVerificationRequests(ctx context.Context) ([]model.VerificationRequest, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetVerificationRequests-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.GetVerificationRequests(nextCtx)
}

func (service *ProfileService) DeleteProfile(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteProfile-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	err := service.ProfileRepository.DeleteProfile(nextCtx, profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	err = service.deleteProfileInAuth(nextCtx, profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	err = service.deleteProfileInConnection(nextCtx, profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	err = service.deleteProfilesPosts(nextCtx, profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	return nil
}

func (service *ProfileService) SendAgentRequest(ctx context.Context, loggedUserID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "SendAgentRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n", loggedUserID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	request := model.AgentRequest{ProfileId: loggedUserID}
	return service.ProfileRepository.SendAgentRequest(nextCtx, &request)
}

func (service *ProfileService) GetAgentRequests(ctx context.Context) ([]dto.ResponseAgentRequestDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAgentRequests-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	requests, err := service.ProfileRepository.GetAgentRequests(nextCtx)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.ResponseAgentRequestDTO, 0)
	for _, value := range requests {
		profile, err := service.ProfileRepository.GetProfileByID(nextCtx, value.ProfileId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		ret = append(ret, dto.ResponseAgentRequestDTO{Username: profile.Username, ProfileID: profile.ID,
			Email: profile.Email, Website: profile.Website})
	}
	return ret, nil
}

func (service *ProfileService) GetProfileUsernamesByIDs(ctx context.Context, ids []string) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileUsernamesByIDs-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	ret := make([]string, 0)
	for _, value := range ids {
		profile, err := service.GetProfileByID(nextCtx, util.String2Uint(value))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				util.Tracer.LogError(span, fmt.Errorf("profile with id %s not found", value))
				ret = append(ret, "")
				continue
			} else {
				util.Tracer.LogError(span, err)
				return nil, err
			}
		}
		ret = append(ret, profile.Username)
	}
	return ret, nil

}

func (service *ProfileService) GetByInterests(ctx context.Context, interests []string) ([]model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetByInterests-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.GetByInterests(nextCtx, interests)
}

func (service *ProfileService) ProcessAgentRequest(ctx context.Context, requestDTO dto.ProcessAgentRequest) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ProcessAgentRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing agent request for profile id %v\n", requestDTO.ProfileID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profileID := util.String2Uint(requestDTO.ProfileID)
	request, err := service.ProfileRepository.GetAgentRequestByProfileID(nextCtx, profileID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, profileID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	subject, message := "Agent Request Rejected", "Your agent request has been declined!"
	if requestDTO.Accept {
		err = makeAgent(nextCtx, request.ProfileId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		subject, message = "Agent Request Accepted", "Your agent request has been accepted!"
	}
	util.Tracer.LogFields(span, "service", fmt.Sprintf("send message with agent request status to mail  %s", profile.Email))
	go util.SendMail(profile.Email, subject, message)
	return service.ProfileRepository.DeleteAgentRequest(nextCtx, request)
}

func (service *ProfileService) GetPersonalDataByProfileId(ctx context.Context, id uint) (*model.PersonalData, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPersonalDataByProfileId-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profile, err := service.ProfileRepository.GetPersonalDataByProfileId(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) GetProfileInterests(ctx context.Context, id uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileInterests-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.GetProfileInterests(nextCtx, id)
}

func (service *ProfileService) deleteProfileInAuth(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "deleteProfileInAuth-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	authHost, authPort := util.GetAuthHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/ban/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (service *ProfileService) deleteProfileInConnection(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "deleteProfileInConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/profile/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (service *ProfileService) deleteProfilesPosts(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "deleteProfilesPosts-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postHost, postPort := util.GetPostHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func (service *ProfileService) Test(ctx context.Context, key string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Test-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing key %s\n", key))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.InsertInRedis(nextCtx, key, "test")
}

func (service *ProfileService) changePrivacyInPostService(ctx context.Context, isPrivate bool, loggedUserId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "changePrivacyInPostService-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v, change is private to %t\n", loggedUserId, isPrivate))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postHost, postPort := util.GetPostHostAndPort()
	type Privacy struct {
		IsPrivate bool `json:"isPrivate"`
	}
	input := Privacy{IsPrivate: isPrivate}
	jsonPrivacy, _ := json.Marshal(input)
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/privacy",
		jsonPrivacy, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
	}
	return err
}

func (service *ProfileService) changeUsernameInPostService(ctx context.Context, loggedUserId uint, username string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "changeUsernameInPostService-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v, change username to %s\n", loggedUserId, username))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postHost, postPort := util.GetPostHostAndPort()
	type UsernameDto struct {
		Username string `json:"username"`
	}
	input := UsernameDto{Username: username}
	jsonUsername, _ := json.Marshal(input)
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/username",
		jsonUsername, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
	}
	return err
}

func (service *ProfileService) GetProfileIdsByUsernames(ctx context.Context, usernames []string) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileIdsByUsernames-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.ProfileRepository.GetProfileIdsByUsernames(nextCtx, usernames)
}

func makeAgent(ctx context.Context, profileID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "makeAgent-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n",profileID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	authHost, authPort := util.GetAuthHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/make-agent/"+util.Uint2String(profileID),
		nil, map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad profile id"))
		return fmt.Errorf("BAD_PROFILE_ID")
	}
	return nil
}
