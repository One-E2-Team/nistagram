package model

import "gorm.io/gorm"

type ProfileSettings struct {
	gorm.Model
	IsPrivate                    bool `json:"isPrivate"`
	CanReceiveMessageFromUnknown bool `json:"CanReceiveMessageFromUnknown"`
	CanBeTagged                  bool `json:"canBeTagged"`
	ProfileID uint
}
