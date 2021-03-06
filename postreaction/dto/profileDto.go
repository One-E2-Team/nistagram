package dto

type ProfileDto struct {
	Username        string          `json:"username"`
	ProfileId       uint            `json:"id"`
	ProfileSettings ProfileSettings `json:"profileSettings"`
}

type ProfileSettings struct {
	IsPrivate    bool `json:"isPrivate"`
	CanBeTagged  bool `json:"canBeTagged"`
}