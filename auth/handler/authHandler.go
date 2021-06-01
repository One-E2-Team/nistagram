package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"nistagram/auth/dto"
	"nistagram/auth/service"
	"time"
)

const ExpiresIn = 86400000
const TokenSecret = "token_secret"

type AuthHandler struct {
	AuthService *service.AuthService
}

type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var dto dto.LogInDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.LogIn(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {

		claims := TokenClaims{Username: dto.Username, StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpiresIn,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth_service",
		}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TokenSecret))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Authorization", signedToken)
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.Register(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Header().Set("Content-Type", "application/json")
}
