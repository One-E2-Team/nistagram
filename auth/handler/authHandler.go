package handler

import (
	"context"
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
	span := util.Tracer.StartSpanFromRequest("LogIn-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var req dto.LogInDTO
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err := json.NewDecoder(r.Body).Decode(&req)
	methodPath := "nistagram/auth/handler.LogIn"
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user *model.User
	user, err = handler.AuthService.LogIn(ctx, req)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.WARN, methodPath, util.GetIPAddress(r), err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := util.CreateToken(user.ProfileId, "auth_service")
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := dto.TokenResponseDTO{
		Token:     token,
		Email:     user.Email,
		ProfileId: user.ProfileId,
		Roles:     user.Roles,
		Username:  user.Username,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJson)
	util.Logging(util.SUCCESS, methodPath, util.GetIPAddress(r), "Successful login for "+util.GetLoggingStringFromID(user.ProfileId), "auth")
	w.Header().Set("Content-Type", "application/json")
}


func (handler *AuthHandler) LogInAgentAPI(writer http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("LogInAgentAPI-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))

	var login dto.APILoginDTO
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		util.Tracer.LogError(span, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	resp, err := handler.AuthService.AgentLoginAPIToken(ctx, login)
	if err != nil {
		util.Tracer.LogError(span, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	respJson, err := json.Marshal(*resp)
	if err != nil {
		util.Tracer.LogError(span, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(respJson)
	writer.Header().Set("Content-Type", "text/plain")
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Register-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var registerDTO dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&registerDTO)
	registerDTO = safeRegisterDTO(registerDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.AuthService.Register(ctx, registerDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) RequestPassRecovery(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("RequestPassRecovery-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var email string
	methodPath := "nistagram/auth/handler.RequestPassRecovery"
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.AuthService.RequestPassRecovery(ctx, email)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.WARN, methodPath, util.GetIPAddress(r), "'"+email+"' "+err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	util.Logging(util.INFO, methodPath, util.GetIPAddress(r), "Successful requested recovery for '"+email+"'", "auth")
	_, _ = w.Write([]byte("Check your email!"))
	w.Header().Set("Content-Type", "text/plain")
}

func (handler *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ChangePassword-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var recoveryDTO dto.RecoveryDTO
	methodPath := "nistagram/auth/handler.ChangePassword"
	err := json.NewDecoder(r.Body).Decode(&recoveryDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.AuthService.ChangePassword(ctx, recoveryDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.WARN, methodPath, util.GetIPAddress(r), err.Error(), "auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	util.Logging(util.INFO, methodPath, util.GetIPAddress(r), "Successful password change for profileId: '"+recoveryDTO.Id+"'", "auth")
	_, _ = w.Write([]byte("Password successfully changed!"))
	w.Header().Set("Content-Type", "text/plain")
}

func (handler *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdateUser-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var updateUserDto dto.UpdateUserDTO
	err := json.NewDecoder(r.Body).Decode(&updateUserDto)
	updateUserDto = safeUpdateUserDto(updateUserDto)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.AuthService.UpdateUser(ctx, updateUserDto)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) ValidateUser(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ValidateUser-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	vars := mux.Vars(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err := handler.AuthService.ValidateUser(ctx, vars["id"], vars["uuid"])
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	frontHost, frontPort := util.GetFrontHostAndPort()
	//_, err = fmt.Fprintf(w, "<html><head><script>window.location.href = \""+util.FrontProtocol+"://"+
	//	frontHost+":"+frontPort+"/web#/2fa-totp/"+vars["qruuid"]+"\";</script></head><body></body></html>")
	_, err = fmt.Fprintf(w, "<html><head><script>window.location.href = \""+util.GetFrontProtocol()+"://"+
		frontHost+":"+frontPort+"/web#/log-in\";</script></head><body></body></html>")
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (handler *AuthHandler) GetAgentAPIToken(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAgentAPIToken-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	apiToken, err := handler.AuthService.GetAgentAPIToken(ctx, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"apiToken\":\"" + apiToken + "\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) GetPrivileges(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPrivileges-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	vars := mux.Vars(r)
	id := util.String2Uint(vars["profileId"])
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	var privileges = handler.AuthService.GetPrivileges(ctx,id)
	if privileges == nil || len(*privileges) == 0 {
		var temp = make([]string, 0)
		//TODO: nzm sta da ispisem
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(temp)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(privileges)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return
		}
	}
}

func (handler *AuthHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("BanUser-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	vars := mux.Vars(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err := handler.AuthService.BanUser(ctx, util.String2Uint(vars["profileID"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AuthHandler) MakeAgent(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("MakeAgent-handler", r)
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	vars := mux.Vars(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err := handler.AuthService.MakeAgent(ctx, util.String2Uint(vars["profileID"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
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
	return template.HTMLEscapeString(str)
}
