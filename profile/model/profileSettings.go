package model

type ProfileSettings struct {
	IsPrivate                    bool `json:"isPrivate"`
	CanReceiveMessageFromUnknown bool `json:"canReceiveMessageFromUnknown"`
	CanBeTagged                  bool `json:"canBeTagged"`
}
