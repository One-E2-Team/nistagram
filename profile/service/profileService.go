package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/profile/dto"
	"nistagram/profile/model"
	"nistagram/profile/repository"
	"nistagram/util"
)

type ProfileService struct {
	ProfileRepository *repository.ProfileRepository
}

func (service *ProfileService) Register(dto dto.RegistrationDto) error {
	profileSettings := model.ProfileSettings{IsPrivate: dto.IsPrivate, CanReceiveMessageFromUnknown: true, CanBeTagged: true}
	personalData := model.PersonalData{Name: dto.Name, Surname: dto.Surname, Telephone: dto.Telephone,
		Gender: dto.Gender, BirthDate: dto.BirthDate}
	for _, item := range dto.InterestedIn {
		interest := service.ProfileRepository.FindInterestByName(item)
		personalData.AddItem(interest)
	}
	profile := model.Profile{Username: dto.Username, Email: dto.Email, ProfileSettings: profileSettings,
		PersonalData: personalData, Biography: dto.Biography, Website: dto.WebSite, IsVerified: false}
	err := service.ProfileRepository.CreateProfile(&profile)
	if err != nil {
		return err
	}
	postBody, _ := json.Marshal(map[string]string{
		"profileId": util.Uint2String(profile.ID),
		"password":  dto.Password,
		"email":     profile.Email,
		"username":  profile.Username,
	})
	//responseBody := bytes.NewBuffer(postBody)
	err = registerInAuth(postBody)
	if err != nil{
		return err
	}
	go func() {
		err := registerInConnection(profile.ID, postBody)
		if err != nil {
			fmt.Println("conn bug")
			fmt.Println(err)
		}
	}()
	return nil
}

func (service *ProfileService) RegisterAgent(dto dto.RegistrationDto) error {
	err := service.Register(dto)
	if err != nil{
		return err
	}

	profile, err := service.GetProfileByUsername(dto.Username)
	if err != nil{
		return err
	}

	fmt.Println("Profile id: ", profile.ID)
	err = makeAgent(profile.ID)

	return err
}

func registerInAuth(postBody []byte) error {
	authHost, authPort := util.GetAuthHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/register", postBody,
		map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func registerInConnection(profileId uint, postBody []byte) error {
	connHost, connPort := util.GetConnectionHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/profile/"+util.Uint2String(profileId), postBody,
		map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (service *ProfileService) Search(username string) []string {
	return service.ProfileRepository.FindUsernameContains(username)
}

func (service *ProfileService) SearchForTag(loggedUserId uint, username string) ([]string, error) {
	var ret []string
	usernames := service.ProfileRepository.FindUsernameContains(username)

	followingProfiles, err := getFollowingProfiles(loggedUserId)
	if err != nil {
		return ret, err
	}

	for i := 0; i < len(usernames); i++ {
		profile, err := service.GetProfileByUsername(usernames[i])
		if err != nil {
			fmt.Println("Can't get profile by username!")
		}
		if util.IsFollowed(followingProfiles, profile.ID) && profile.ProfileSettings.CanBeTagged {
			ret = append(ret, profile.Username)
		}
	}

	return ret, nil
}

func getFollowingProfiles(loggedUserId uint) ([]util.FollowingProfileDTO, error) {
	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})

	if err != nil {
		return nil, err
	}

	var followingProfiles []util.FollowingProfileDTO

	err = json.NewDecoder(resp.Body).Decode(&followingProfiles)
	if err != nil{
		return nil, err
	}

	return followingProfiles, err
}

func (service *ProfileService) GetProfileByUsername(username string) (*model.Profile, error) {
	profile, err := service.ProfileRepository.FindProfileByUsername(username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) ChangeProfileSettings(dto dto.ProfileSettingsDTO, loggedUserId uint) error {
	profile, err := service.ProfileRepository.GetProfileByID(context.Background(), loggedUserId)
	if err != nil {
		return err
	}
	profileSettings := profile.ProfileSettings
	if profileSettings.IsPrivate != dto.IsPrivate {
		err = service.changePrivacyInPostService(dto.IsPrivate, loggedUserId)
		if err != nil {
			return err
		}
	}
	profileSettings.IsPrivate = dto.IsPrivate
	profileSettings.CanBeTagged = dto.CanBeTagged
	profileSettings.CanReceiveMessageFromUnknown = dto.CanReceiveMessageFromUnknown
	err = service.ProfileRepository.UpdateProfileSettings(profileSettings)
	return err
}

func (service *ProfileService) ChangePersonalData(dto dto.PersonalDataDTO, loggedUserId uint) (string, string, error) {
	profile, err := service.ProfileRepository.GetProfileByID(context.Background(), loggedUserId)
	if err != nil {
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
	err = service.ProfileRepository.UpdateProfile(profile)
	if err != nil {
		return "", "", err
	}
	if callAuth {
		postBody, _ := json.Marshal(map[string]string{
			"profileId": util.Uint2String(profile.ID),
			"email":     profile.Email,
			"username":  profile.Username,
		})
		authHost, authPort := util.GetAuthHostAndPort()
		_, err = util.CrossServiceRequest(http.MethodPost,
			util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/update-user", postBody,
			map[string]string{"Content-Type": "application/json;"})
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
	}
	if callPost {
		err = service.changeUsernameInPostService(loggedUserId, dto.Username)
		if err != nil {
			return "", "", err
		}
	}
	err = service.ProfileRepository.UpdatePersonalData(profile.PersonalData)
	return oldUsername, oldEmail, err
}

func (service *ProfileService) GetAllInterests() ([]string, error) {
	interests, err := service.ProfileRepository.GetAllInterests()
	return interests, err
}

func (service *ProfileService) GetAllCategories() ([]string, error) {
	categories, err := service.ProfileRepository.GetAllCategories()
	return categories, err
}

func (service *ProfileService) CreateVerificationRequest(profileId uint, requestDTO dto.VerificationRequestDTO, fileName string) error {
	category, err := service.ProfileRepository.GetCategoryByName(requestDTO.Category)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var verReq = model.VerificationRequest{ProfileID: profileId, Name: requestDTO.Name, Surname: requestDTO.Surname,
		VerificationStatus: model.SENT, ImagePath: fileName, Category: *category}
	err = service.ProfileRepository.CreateVerificationRequest(&verReq)
	return err
}

func (service *ProfileService) UpdateVerificationRequest(verifyDTO dto.VerifyDTO) error {
	request, err := service.ProfileRepository.GetVerificationRequestById(verifyDTO.VerificationId)
	if err != nil {
		return err
	}
	if verifyDTO.Status {
		request.VerificationStatus = model.VERIFIED
		err = service.ProfileRepository.UpdateVerificationRequest(*request)
		if err != nil {
			return err
		}
		profile, err := service.ProfileRepository.GetProfileByID(context.Background(), request.ProfileID)
		if err != nil {
			return err
		}
		profile.IsVerified = true
		err = service.ProfileRepository.UpdateProfile(profile)
	} else {
		err = service.ProfileRepository.DeleteVerificationRequest(request)
	}

	return err
}

func (service *ProfileService) GetMyProfileSettings(loggedUserId uint) (dto.ProfileSettingsDTO, error) {
	ret := dto.ProfileSettingsDTO{}
	profile, err := service.ProfileRepository.GetProfileByID(context.Background(), loggedUserId)
	if err != nil {
		return ret, err
	}
	ret.CanReceiveMessageFromUnknown = profile.ProfileSettings.CanReceiveMessageFromUnknown
	ret.CanBeTagged = profile.ProfileSettings.CanBeTagged
	ret.IsPrivate = profile.ProfileSettings.IsPrivate
	return ret, nil
}

func (service *ProfileService) GetMyPersonalData(loggedUserId uint) (dto.PersonalDataDTO, error) {
	profile, err := service.ProfileRepository.GetProfileByID(context.Background(), loggedUserId)
	if err != nil {
		return dto.PersonalDataDTO{}, err
	}
	personalData := profile.PersonalData
	ret := dto.PersonalDataDTO{Username: profile.Username, Name: personalData.Name, Surname: personalData.Surname,
		Email: profile.Email, Telephone: personalData.Telephone, Gender: personalData.Gender,
		BirthDate: personalData.BirthDate, Biography: personalData.BirthDate, Website: profile.Website}
	return ret, nil
}

func (service *ProfileService) GetProfileByID(ctx context.Context, id uint) (*model.Profile, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfilByID-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("sevicing get profile for id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	profile, err := service.ProfileRepository.GetProfileByID(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) GetVerificationRequests() ([]model.VerificationRequest, error) {
	return service.ProfileRepository.GetVerificationRequests()
}

func (service *ProfileService) DeleteProfile(profileId uint) error {
	err := service.ProfileRepository.DeleteProfile(profileId)
	if err != nil {
		return err
	}

	err = service.deleteProfileInAuth(profileId)
	if err != nil {
		return err
	}

	err = service.deleteProfileInConnection(profileId)
	if err != nil {
		return err
	}

	err = service.deleteProfilesPosts(profileId)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProfileService) SendAgentRequest(loggedUserID uint) error {
	request := model.AgentRequest{ProfileId: loggedUserID}
	return service.ProfileRepository.SendAgentRequest(&request)
}

func (service *ProfileService) GetAgentRequests() ([]dto.ResponseAgentRequestDTO, error) {
	requests, err := service.ProfileRepository.GetAgentRequests()
	if err != nil {
		return nil, err
	}
	ret := make([]dto.ResponseAgentRequestDTO, 0)
	for _, value := range requests {
		profile, err := service.ProfileRepository.GetProfileByID(context.Background(), value.ProfileId)
		if err != nil {
			return nil, err
		}
		ret = append(ret, dto.ResponseAgentRequestDTO{Username: profile.Username, ProfileID: profile.ID,
			Email: profile.Email, Website: profile.Website})
	}
	return ret, nil
}

func (service *ProfileService) GetProfileUsernamesByIDs(ids []string) ([]string, error) {
	ret := make([]string, 0)
	for _, value := range ids {
		profile, err := service.GetProfileByID(context.Background(), util.String2Uint(value))
		if err != nil {
			return nil, err
		}
		ret = append(ret, profile.Username)
	}
	return ret, nil

}

func (service *ProfileService) GetByInterests(interests []string) ([]model.Profile, error) {
	return service.ProfileRepository.GetByInterests(interests)
}

func (service *ProfileService) ProcessAgentRequest(requestDTO dto.ProcessAgentRequest) error {
	profileID := util.String2Uint(requestDTO.ProfileID)
	request, err := service.ProfileRepository.GetAgentRequestByProfileID(profileID)
	if err != nil {
		return err
	}
	profile, err := service.ProfileRepository.GetProfileByID(context.Background(), profileID)
	if err != nil {
		return err
	}
	subject, message := "Agent Request Rejected", "Your agent request has been declined!"
	if requestDTO.Accept {
		err = makeAgent(request.ProfileId)
		if err != nil {
			return err
		}
		subject, message = "Agent Request Accepted", "Your agent request has been accepted!"
	}
	go util.SendMail(profile.Email, subject, message)
	return service.ProfileRepository.DeleteAgentRequest(request)
}

func (service *ProfileService) GetPersonalDataByProfileId(id uint) (*model.PersonalData, error) {
	profile, err := service.ProfileRepository.GetPersonalDataByProfileId(id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) deleteProfileInAuth(profileId uint) error {
	authHost, authPort := util.GetAuthHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/ban/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func (service *ProfileService) deleteProfileInConnection(profileId uint) error {
	connHost, connPort := util.GetConnectionHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/profile/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func (service *ProfileService) deleteProfilesPosts(profileId uint) error {
	postHost, postPort := util.GetPostHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(profileId), nil,
		map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func (service *ProfileService) Test(key string) error {
	return service.ProfileRepository.InsertInRedis(key, "test")
}

func (service *ProfileService) changePrivacyInPostService(isPrivate bool, loggedUserId uint) error {
	postHost, postPort := util.GetPostHostAndPort()
	type Privacy struct {
		IsPrivate bool `json:"isPrivate"`
	}
	input := Privacy{IsPrivate: isPrivate}
	jsonPrivacy, _ := json.Marshal(input)
	_, err := util.CrossServiceRequest(http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/privacy",
		jsonPrivacy, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (service *ProfileService) changeUsernameInPostService(loggedUserId uint, username string) error {
	postHost, postPort := util.GetPostHostAndPort()
	type UsernameDto struct {
		Username string `json:"username"`
	}
	input := UsernameDto{Username: username}
	jsonUsername, _ := json.Marshal(input)
	_, err := util.CrossServiceRequest(http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/username",
		jsonUsername, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (service *ProfileService) GetProfileIdsByUsernames(usernames []string) ([]string, error) {
	return service.ProfileRepository.GetProfileIdsByUsernames(usernames)
}

func makeAgent(profileID uint) error {
	authHost, authPort := util.GetAuthHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodPut,
		util.GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/make-agent/"+util.Uint2String(profileID),
		nil, map[string]string{})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("BAD_PROFILE_ID")
	}
	return nil
}
