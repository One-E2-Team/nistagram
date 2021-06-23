package dto

import "nistagram/agent/model"

type TokenResponseDTO struct {
	Token     string       `json:"token"`
	Email     string       `json:"email"`
	UserId 	  uint         `json:"userId"`
	Roles     []model.Role `json:"roles"`
}
