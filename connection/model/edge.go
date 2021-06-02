package model

type Connection struct {
	PrimaryProfile		Profile	`json:"primary"`
	SecondaryProfile	Profile	`json:"secondary"`
	Muted				bool	`json:"muted"`
	CloseFriend			bool	`json:"closeFriend"`
	NotifyPost			bool	`json:"notifyPost"`
	NotifyStory			bool	`json:"notifyStory"`
	NotifyMessage		bool	`json:"notifyMessage"`
	NotifyComment		bool	`json:"notifyComment"`
	Approved			bool	`json:"approved"`
	MessageConnected	bool	`json:"messageConnected"`
	Block				bool	`json:"block"`
}
