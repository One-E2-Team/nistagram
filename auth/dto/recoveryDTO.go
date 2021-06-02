package dto

type RecoveryDTO struct {
	Id       uint   `json:"id"`
	Uuid     string `json:"uuid"`
	Password string `json:"password"`
}
