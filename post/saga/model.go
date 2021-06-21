package saga

import (
	"encoding/json"
	"nistagram/profile/model"
)

const (
	ProfileChannel string = "ProfileChannel"
	AuthChannel string = "AuthChannel"
	ConnectionChannel string = "ConnectionChannel"
	PostChannel string = "PostChannel"
	ReplyChannel    string = "ReplyChannel"
	ActionStart     string = "Start"
	ActionRollback  string = "RollbackMsg"
	ActionDone      string = "DoneMsg"
	ActionError     string = "ErrorMsg"
	ChangeProfilesPrivacy string = "ChangeProfilesPrivacy"
	RegisterProfile string = "RegisterProfile"
	DeleteProfile string = "DeleteProfile"
	ProfileService string = "ProfileService"
	AuthService string = "AuthService"
	ConnectionService string = "ConnectionService"
	PostService string = "PostService"
)

type Message struct {
	NextService		string		    `json:"nextService"`
	SenderService   string          `json:"sender_service"`
	Action          string          `json:"action"`
	Functionality   string		    `json:"functionality"`
	Profile			model.Profile   `json:"profile"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
