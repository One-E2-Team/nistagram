package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"nistagram/auth/dto"
	"nistagram/auth/service"
	"strings"
	"time"
)

const ExpiresIn = 86400
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
			ExpiresAt: time.Now().Unix() + ExpiresIn,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth_service",
		}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TokenSecret))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Authorization", "Bearer "+signedToken)
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

func (handler *AuthHandler) GetToken(header http.Header) (string, error) {
	reqToken := header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("NO_TOKEN")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

func (handler *AuthHandler) GetLoggedUsername(r *http.Request) (string, error) {
	tokenString, err := handler.GetToken(r.Header)
	if err != nil {
		return "", fmt.Errorf("NO_TOKEN")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecret), nil
	})
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.Username, nil
	} else {
		fmt.Println(err.Error())
		return "", nil
	}
}
