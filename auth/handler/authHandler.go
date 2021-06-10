package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
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
		ProfileId: user.ProfileId,
		Roles: user.Roles,
		Username: user.Username,
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
	registerDTO = safeRegisterDTO(registerDTO)
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
	w.Header().Set("Content-Type", "text/plain")
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
	w.Header().Set("Content-Type", "text/plain")
}

func (handler *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserDto dto.UpdateUserDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserDto)
	updateUserDto = safeUpdateUserDto(updateUserDto)
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
	//TODO: parametrize host
	_, err = fmt.Fprintf(w, "<html><head><script>window.location.href = \"http://localhost:81/web#/log-in\";</script></head><body></body></html>")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (handler *AuthHandler) GetPrivileges(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := util.String2Uint(vars["profileId"])
	var privileges *[]string = handler.AuthService.GetPrivileges(id)
	if privileges == nil  || len(*privileges) == 0 {
		var temp []string = make([]string, 0)
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(temp)
	} else {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(privileges)
	}
}


func safeRegisterDTO(dto dto.RegisterDTO) dto.RegisterDTO {
	dto.Username = sanitize(dto.Username)
	dto.Email = sanitize(dto.Email)
	return dto
}

func safeUpdateUserDto(dto dto.UpdateUserDTO) dto.UpdateUserDTO {
	dto.Username = sanitize(dto.Username)
	dto.Email = sanitize(dto.Email)
	return dto
}

func sanitize(str string) string {
	return  template.HTMLEscapeString(str)
}
