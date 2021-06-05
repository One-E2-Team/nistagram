package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/service"
	"nistagram/util"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var req dto.LogInDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	var user *model.User
	user, err = handler.AuthService.LogIn(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	token, err := util.CreateToken(user.ProfileId, "auth_service")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	resp := dto.TokenResponseDTO{
		Token: token,
		Email: user.Email,
		Roles: user.Roles,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJson)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	err = handler.AuthService.Register(registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) RequestPassRecovery(w http.ResponseWriter, r *http.Request) {
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	err = handler.AuthService.RequestPassRecovery(email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Check your email!"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var recoveryDTO dto.RecoveryDTO
	err := json.NewDecoder(r.Body).Decode(&recoveryDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	err = handler.AuthService.ChangePassword(recoveryDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Password successfully changed!"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserDto dto.UpdateUserDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	err = handler.AuthService.UpdateUser(updateUserDto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) ValidateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := handler.AuthService.ValidateUser(vars["id"], vars["uuid"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}
