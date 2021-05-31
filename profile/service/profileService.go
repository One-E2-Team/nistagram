package service

import (
	"nistagram/profile/dto"
	"nistagram/profile/model"
	"nistagram/profile/repository"
)

type ProfileService struct {
	ProfileRepository *repository.ProfileRepository
}

func (service *ProfileService) Register(dto dto.RegistrationDto) error{
	profileSettings := model.ProfileSettings{IsPrivate: dto.IsPrivate, CanRecieveMessageFromUnknown: true, CanBeTagged: true}
	personalData := model.PersonalData{Name: dto.Name, Surname: dto.Surname,
		Email: dto.Email, Telephone: dto.Telephone, Gender: dto.Gender, BirthDate: dto.BirthDate}
	for _, item := range dto.InterestedIn{
		interest := service.ProfileRepository.FindInterestByName(item)
		personalData.AddItem(interest)
	}
	profile := model.Profile{Username: dto.Username,ProfileSettings: profileSettings,PersonalData: personalData, Biography: dto.Biography, Website: dto.WebSite, Type: model.REGULAR}
	err := service.ProfileRepository.CreateProfile(&profile)
	if err != nil{
		return err
	}
	return nil
}

