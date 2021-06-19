package model

type ConnectionEdge struct {
	PrimaryProfile    uint `json:"primary"`
	SecondaryProfile  uint `json:"secondary"`
	Muted             bool `json:"muted"`
	CloseFriend       bool `json:"closeFriend"`
	NotifyPost        bool `json:"notifyPost"`
	NotifyStory       bool `json:"notifyStory"`
	NotifyComment     bool `json:"notifyComment"`
	ConnectionRequest bool `json:"connectionRequest"`
	Approved          bool `json:"approved"`
}

func (conn *ConnectionEdge) ToMap() map[string]interface{} {
	return ToMap(conn)
}

type BlockEdge struct {
	PrimaryProfile    uint `json:"primary"`
	SecondaryProfile  uint `json:"secondary"`
}

func (block *BlockEdge) ToMap() map[string]interface{} {
	return ToMap(block)
}

type MessageEdge struct {
	PrimaryProfile		uint `json:"primary"`
	SecondaryProfile	uint `json:"secondary"`
	Approved			bool `json:"approved"`
	NotifyMessage		bool `json:"notifyMessage"`
}

func (message *MessageEdge) ToMap() map[string]interface{} {
	return ToMap(message)
}