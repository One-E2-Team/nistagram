package service

import (
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
		PersonalData: personalData, Biography: dto.Biography, Website: dto.WebSite, Type: model.REGULAR}
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
	go func() {
		err := registerInAuth(postBody)
		if err != nil {
			fmt.Println("auth bug")
			fmt.Println(err)
		}
	}()
	go func() {
		err := registerInConnection(profile.ID, postBody)
		if err != nil {
			fmt.Println("conn bug")
			fmt.Println(err)
		}
	}()
	return nil
}

func registerInAuth(postBody []byte) error {
	authHost, authPort := util.GetAuthHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodPost,
		util.CrossServiceProtocol+"://"+authHost+":"+authPort+"/register", postBody,
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
		util.CrossServiceProtocol+"://"+connHost+":"+connPort+"/profile/"+util.Uint2String(profileId), postBody,
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

func (service *ProfileService) GetProfileByUsername(username string) (*model.Profile, error) {
	profile, err := service.ProfileRepository.FindProfileByUsername(username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) ChangeProfileSettings(dto dto.ProfileSettingsDTO, loggedUserId uint) error {
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
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
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
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
			util.CrossServiceProtocol+"://"+authHost+":"+authPort+"/update-user", postBody,
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

func (service *ProfileService) GetMyProfileSettings(loggedUserId uint) (dto.ProfileSettingsDTO, error) {
	ret := dto.ProfileSettingsDTO{}
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
	if err != nil {
		return ret, err
	}
	ret.CanReceiveMessageFromUnknown = profile.ProfileSettings.CanReceiveMessageFromUnknown
	ret.CanBeTagged = profile.ProfileSettings.CanBeTagged
	ret.IsPrivate = profile.ProfileSettings.IsPrivate
	return ret, nil
}

func (service *ProfileService) GetMyPersonalData(loggedUserId uint) (dto.PersonalDataDTO, error) {
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
	if err != nil {
		return dto.PersonalDataDTO{}, err
	}
	personalData := profile.PersonalData
	ret := dto.PersonalDataDTO{Username: profile.Username, Name: personalData.Name, Surname: personalData.Surname,
		Email: profile.Email, Telephone: personalData.Telephone, Gender: personalData.Gender,
		BirthDate: personalData.BirthDate, Biography: personalData.BirthDate, Website: profile.Website}
	return ret, nil
}

func (service *ProfileService) GetProfileByID(id uint) (*model.Profile, error) {
	profile, err := service.ProfileRepository.GetProfileByID(id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (service *ProfileService) GetVerificationRequests() ([]model.VerificationRequest, error) {
	return service.ProfileRepository.GetVerificationRequests()
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/privacy",
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/user/"+util.Uint2String(loggedUserId)+"/username",
		jsonUsername, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		fmt.Println(err)
	}
	return err
}
