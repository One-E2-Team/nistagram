package model

import "gorm.io/gorm"

type ProfileSettings struct {
	gorm.Model
	IsPrivate                    bool `json:"isPrivate"`
	CanReceiveMessageFromUnknown bool `json:"canReceiveMessageFromUnknown"`
	CanBeTagged                  bool `json:"canBeTagged"`
	ProfileID                    uint
}
