package service

import (
	"bytes"
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
	profile := model.Profile{Username: dto.Username, Email: dto.Email, ProfileSettings: profileSettings, PersonalData: personalData, Biography: dto.Biography, Website: dto.WebSite, Type: model.REGULAR}
	err := service.ProfileRepository.CreateProfile(&profile)
	if err != nil {
		return err
	}
	postBody, _ := json.Marshal(map[string]string{
		"profileId": util.Uint2String(profile.ID),
		"password":  dto.Password,
		"email":     profile.Email,
	})
	responseBody := bytes.NewBuffer(postBody)
	authHost, authPort := util.GetAuthHostAndPort()
	_, err = http.Post("http://"+authHost+":"+authPort+"/register", "application/json", responseBody)
	if err != nil {
		fmt.Println(err)
		return err
	}
	connHost, connPort := util.GetConnectionHostAndPort()
	_, err = http.Post("http://"+connHost+":"+connPort+"/profile/" + util.Uint2String(profile.ID), "application/json", responseBody)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (service *ProfileService) Search(username string) []string {
	return service.ProfileRepository.FindUsernameContains(username)
}

func (service *ProfileService) GetProfileByUsername(username string) *model.Profile {
	return service.ProfileRepository.FindProfileByUsername(username)
}

func (service *ProfileService) ChangeProfileSettings(dto dto.ProfileSettingsDTO, loggedUserId uint) error {
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
	if err != nil {
		return err
	}
	profileSettings := profile.ProfileSettings
	if profileSettings.IsPrivate != dto.IsPrivate {
		//TODO: change data in post ms
	}
	profileSettings.IsPrivate = dto.IsPrivate
	profileSettings.CanBeTagged = dto.CanBeTagged
	profileSettings.CanReceiveMessageFromUnknown = dto.CanReceiveMessageFromUnknown
	service.ProfileRepository.UpdateProfileSettings(profileSettings)
	return nil
}

func (service *ProfileService) ChangePersonalData(dto dto.PersonalDataDTO, loggedUserId uint) error {
	profile, err := service.ProfileRepository.GetProfileByID(loggedUserId)
	if err != nil {
		return err
	}
	callAuth := bool(false)
	if profile.Email != dto.Email {
		callAuth = true
	}
	if profile.Username != dto.Username {
		//TODO: change data in other ms
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
		return err
	}
	if callAuth {
		postBody, _ := json.Marshal(map[string]string{
			"profileId": util.Uint2String(profile.ID),
			"email":     profile.Email,
		})
		responseBody := bytes.NewBuffer(postBody)
		authHost, authPort := util.GetAuthHostAndPort()
		_, err = http.Post("http://"+authHost+":"+authPort+"/update-user", "application/json", responseBody)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	err = service.ProfileRepository.UpdatePersonalData(profile.PersonalData)
	return err
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

func (service *ProfileService) Test(key string) error {
	return service.ProfileRepository.InsertInRedis(key, "test")
}
