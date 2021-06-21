package saga

import (
	"encoding/json"
	"nistagram/profile/model"
)

type Message struct {
	NextService		string		    `json:"nextService"`
	SenderService   string          `json:"sender_service"`
	Action          string          `json:"action"`
	Functionality   string		    `json:"functionality"`
	Profile			*model.Profile  `json:"profile"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}