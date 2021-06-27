package dto

type FollowingProfileDTO struct {
	ProfileID        uint     `json:"profileID"`
	CloseFriend 	 bool     `json:"closeFriend"`
}