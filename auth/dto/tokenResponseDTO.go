package dto

import "nistagram/auth/model"

type TokenResponseDTO struct {
	Token     string       `json:"token"`
	Email     string       `json:"email"`
	Roles     []model.Role `json:"roles"`
}
