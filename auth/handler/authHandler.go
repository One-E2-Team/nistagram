package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"nistagram/auth/dto"
	"nistagram/auth/model"
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
	Username     string `json:"username"`
	LoggedUserId uint   `json:"loggedUserId"`
	jwt.StandardClaims
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var dto dto.LogInDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user *model.User
	user, err = handler.AuthService.LogIn(dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := TokenClaims{Username: dto.Username, LoggedUserId: user.ID, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + ExpiresIn,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth_service",
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Authorization", "Bearer "+signedToken)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.Register(dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) RequestPassRecovery(w http.ResponseWriter, r *http.Request) {
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.RequestPassRecovery(username)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Check your email!"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var dto dto.RecoveryDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.ChangePassword(dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password successfully changed!"))
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

func (handler *AuthHandler) GetLoggedUserID(r *http.Request) (uint, error) {
	tokenString, err := handler.GetToken(r.Header)
	if err != nil {
		return 0, fmt.Errorf("NO_TOKEN")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecret), nil
	})
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.LoggedUserId, nil
	} else {
		fmt.Println(err)
		return 0, fmt.Errorf("NO_TOKEN")
	}
}
