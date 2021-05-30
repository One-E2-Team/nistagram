package model

import "gorm.io/gorm"

type ProfileSettings struct {
	gorm.Model
	IsPrivate bool `json:"isPrivate"`
	CanRecieveMessageFromUnknown bool `json:"canRecieveMessageFromUnknown"`
	CanBeTagged bool `json:"vanBeTagged"`
}
