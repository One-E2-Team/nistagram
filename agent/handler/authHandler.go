package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/service"
	"nistagram/agent/util"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.Register(registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var req dto.LogInDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user *model.User
	user, err = handler.AuthService.LogIn(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := util.CreateToken(user.ID, "agent_app")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := dto.TokenResponseDTO{
		Token:     token,
		Email:     user.Email,
		UserId:    user.ID,
		Roles:     user.Roles,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJson)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) ValidateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := handler.AuthService.ValidateUser(vars["id"], vars["uuid"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	frontHost, frontPort := util.GetFrontHostAndPort()
	_, err = fmt.Fprintf(w, "<html><head><script>window.location.href = \""+util.GetFrontProtocol()+"://"+
		frontHost+":"+frontPort+"/#/log-in\";</script></head><body></body></html>")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (handler *AuthHandler) CreateAPIToken(w http.ResponseWriter, r *http.Request) {
	var apiToken string
	err := json.NewDecoder(r.Body).Decode(&apiToken)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.CreateAPIToken(apiToken, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}