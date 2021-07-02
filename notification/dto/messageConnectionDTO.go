package dto

type MessageConnectionDTO struct{
	ProfileId			uint		`json:"profileId"`
	Username			string		`json:"username"`
	MessageApproved		bool		`json:"messageApproved"`
}
