package dto

type WSRequestBodyDTO struct {
	Jwt		string	`json:"jwt"`
	Request string	`json:"request"`
	Data 	string	`json:"data"`
}