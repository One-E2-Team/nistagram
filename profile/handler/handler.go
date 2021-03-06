package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"nistagram/profile/dto"
	"nistagram/profile/service"
	"nistagram/util"
	"os"
	"regexp"
	"strings"
)

type Handler struct {
	ProfileService *service.ProfileService
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Register-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var req dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while decoding.\"}"))
		return
	}

	v := validator.New()

	handler.checkCommonPass(ctx, w, v)

	handler.checkWeakPass(ctx, v, err)

	handler.checkUsername(ctx, v)

	err = v.Struct(req)

	if err != nil {
		util.Tracer.LogError(span, fmt.Errorf("invalid data"))
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Invalid data.\",\n"))
		_, _ = w.Write([]byte("\"errors\":\""))
		for _, e := range err.(validator.ValidationErrors) {
			_, _ = w.Write([]byte(e.Field()))
			_, _ = w.Write([]byte(" "))
		}
		_, _ = w.Write([]byte("\"}"))
		return
	}

	req = safeRegistrationDto(req)
	err = handler.ProfileService.Register(ctx, req)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while registering.\"}"))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}
}

func (handler *Handler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("RegisterAgent-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var req dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while decoding.\"}"))
		return
	}

	v := validator.New()

	handler.checkCommonPass(ctx, w, v)

	handler.checkWeakPass(ctx, v, err)

	handler.checkUsername(ctx, v)

	err = v.Struct(req)

	if err != nil {
		util.Tracer.LogError(span, fmt.Errorf("invalid data"))
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Invalid data.\",\n"))
		_, _ = w.Write([]byte("\"errors\":\""))
		for _, e := range err.(validator.ValidationErrors) {
			_, _ = w.Write([]byte(e.Field()))
			_, _ = w.Write([]byte(" "))
		}
		_, _ = w.Write([]byte("\"}"))
		return
	}

	req = safeRegistrationDto(req)
	err = handler.ProfileService.RegisterAgent(ctx, req)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while registering.\"}"))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}
}

func (handler *Handler) checkUsername(ctx context.Context, v *validator.Validate) {
	span := util.Tracer.StartSpanFromContext(ctx, "checkUsername-handler")
	defer util.Tracer.FinishSpan(span)
	_ = v.RegisterValidation("bad_username", func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 3 || len(fl.Field().String()) > 15 {
			util.Tracer.LogError(span, fmt.Errorf("not valide username length"))
			return false
		}
		ret, _ := regexp.MatchString("([*!@#$%^&(){}\\[:;\\]<>,.?~+\\-\\\\=|/ ])", fl.Field().String())
		if ret {
			util.Tracer.LogError(span, fmt.Errorf("username contains illegal caracters"))
			return false
		}
		return true
	})
}

func (handler *Handler) checkWeakPass(ctx context.Context, v *validator.Validate, err error) {
	span := util.Tracer.StartSpanFromContext(ctx, "checkWeakPass-handler")
	defer util.Tracer.FinishSpan(span)
	_ = v.RegisterValidation("weak_pass", func(fl validator.FieldLevel) bool {
		//^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/])[A-Za-z0-9*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/]{8,}$
		if len(fl.Field().String()) < 8 {
			util.Tracer.LogError(span, fmt.Errorf("password length is less than 8"))
			return false
		}
		ret, _ := regexp.MatchString("(.*[a-z].*)", fl.Field().String())
		if ret == false {
			util.Tracer.LogError(span, fmt.Errorf("password doesn't contain lowercase letters  "))
			return false
		}
		ret, _ = regexp.MatchString("(.*[A-Z].*)", fl.Field().String())
		if ret == false {
			util.Tracer.LogError(span, fmt.Errorf("password doesn't contain uppercase letters  "))
			return false
		}
		ret, _ = regexp.MatchString("(.*[0-9].*)", fl.Field().String())
		if ret == false {
			util.Tracer.LogError(span, fmt.Errorf("password doesn't contain numbers  "))
			return false
		}
		ret, _ = regexp.MatchString("(.*[*!@#$%^&(){}\\[:;\\]<>,.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if err != nil {
			util.Tracer.LogError(span, fmt.Errorf("password doesn't contain special characters"))
			fmt.Println(err)
			return false
		}
		return ret
	})
}

func (handler *Handler) checkCommonPass(ctx context.Context, w http.ResponseWriter, v *validator.Validate) {
	span := util.Tracer.StartSpanFromContext(ctx, "checkCommonPass-handler")
	defer util.Tracer.FinishSpan(span)
	_ = v.RegisterValidation("common_pass", func(fl validator.FieldLevel) bool {
		f, err := os.Open("../common_pass.txt")
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			w.WriteHeader(http.StatusOK)
			return false
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if strings.Contains(fl.Field().String(), scanner.Text()) {
				util.Tracer.LogError(span, fmt.Errorf("common password error"))
				return false
			}
		}
		return true
	})
}

func (handler *Handler) CreateVerificationRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("CreateVerificationRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	profileId := util.GetLoggedUserIDFromToken(r)
	err := r.ParseMultipartForm(0)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	var reqDto dto.VerificationRequestDTO
	data := r.MultipartForm.Value["data"]
	err = json.Unmarshal([]byte(data[0]), &reqDto)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	picture, picHeader, err := r.FormFile("picture")
	uid := uuid.NewString()

	fileSplitted := strings.Split(picHeader.Filename, ".")
	fileName := uid + "." + fileSplitted[1]

	err = handler.ProfileService.CreateVerificationRequest(ctx, profileId, reqDto, fileName)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	f, err := os.OpenFile("../../nistagramstaticdata/data/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	defer func(picture multipart.File) {
		_ = picture.Close()
	}(picture)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	_, err = io.Copy(f, picture)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
}

func (handler *Handler) Search(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Search-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	vars := mux.Vars(r)
	result := handler.ProfileService.Search(ctx, template.HTMLEscapeString(vars["username"]))
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) SearchForTag(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SearchForTag-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	vars := mux.Vars(r)
	loggedUserId := util.GetLoggedUserIDFromToken(r)
	result, err := handler.ProfileService.SearchForTag(ctx, loggedUserId, template.HTMLEscapeString(vars["username"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfileByUsername-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	vars := mux.Vars(r)
	result, err := handler.ProfileService.GetProfileByUsername(ctx, template.HTMLEscapeString(vars["username"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) ChangeProfileSettings(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ChangeProfileSettings-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var req dto.ProfileSettingsDTO
	methodPath := "nistagram/profile/handler.ChangeProfileSettings"
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	err = handler.ProfileService.ChangeProfileSettings(ctx, req, userId)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	util.Logging(util.INFO, methodPath, util.GetIPAddress(r), "Successful profile settings changed, "+util.GetLoggingStringFromID(userId), "profile")
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) UpdateVerificationRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdateVerificationRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var req dto.VerifyDTO
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return
	}

	err = handler.ProfileService.UpdateVerificationRequest(ctx, req)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return
	}

	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
}

func (handler *Handler) ChangePersonalData(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ChangePersonalData-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var req dto.PersonalDataDTO
	methodPath := "nistagram/profile/handler.ChangePersonalData"
	err := json.NewDecoder(r.Body).Decode(&req)
	req = safePersonalDataDto(req)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	oldUsername, oldEmail, err := handler.ProfileService.ChangePersonalData(ctx, req, userId)
	if err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	util.Logging(util.INFO, methodPath, util.GetIPAddress(r), "Successful personal data changed, "+util.GetLoggingStringFromID(userId), "profile")
	if oldUsername != "" {
		util.Logging(util.INFO, methodPath, util.GetIPAddress(r), util.GetLoggingStringFromID(userId)+
			" changed username: "+oldUsername+"->"+req.Username, "profile")
	}
	if oldEmail != "" {
		util.Logging(util.INFO, methodPath, util.GetIPAddress(r), util.GetLoggingStringFromID(userId)+
			" changed email: "+oldEmail+"->"+req.Email, "profile")
	}
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) Test(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Test-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var key string
	err := json.NewDecoder(r.Body).Decode(&key)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProfileService.Test(ctx, key)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetAllInterests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllInterests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	interests, err := handler.ProfileService.GetAllInterests(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(interests)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllCategories-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	categories, err := handler.ProfileService.GetAllCategories(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) GetMyProfileSettings(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyProfileSettings-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	settings, err := handler.ProfileService.GetMyProfileSettings(ctx, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}

}

func (handler *Handler) GetMyPersonalData(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyPersonalData-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	data, err := handler.ProfileService.GetMyPersonalData(ctx, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfileByID-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	var id uint
	vars := mux.Vars(r)
	id = util.String2Uint(vars["id"])
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	profile, err := handler.ProfileService.GetProfileByID(ctx, id)
	fmt.Println(*profile)
	if err != nil {
		fmt.Println(err)
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(*profile)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) GetVerificationRequests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetVerificationRequests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	requests, err := handler.ProfileService.GetVerificationRequests(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
	fmt.Println(requests, err)
}

func (handler *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteProfile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	profileId := util.String2Uint(vars["id"])
	err := handler.ProfileService.DeleteProfile(ctx, profileId)

	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SendAgentRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SendAgentRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	loggedUserID := util.GetLoggedUserIDFromToken(r)
	err := handler.ProfileService.SendAgentRequest(ctx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetAgentRequests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAgentRequests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	requests, err := handler.ProfileService.GetAgentRequests(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(requests)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetProfileUsernamesByIDs(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfileUsernamesByIDs-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	type data struct {
		Ids []string `json:"ids"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret, err := handler.ProfileService.GetProfileUsernamesByIDs(ctx, input.Ids)
	fmt.Println(ret)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(ret)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetByInterests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetByInterests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	type data struct {
		Interests []string `json:"interests"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profiles, err := handler.ProfileService.GetByInterests(ctx, input.Interests)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(profiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) ProcessAgentRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ProcessAgentRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var input dto.ProcessAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := handler.ProfileService.ProcessAgentRequest(ctx, input)
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

func (handler *Handler) GetProfileIdsByUsernames(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfileIdsByUsernames-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	type data struct {
		Usernames []string `json:"usernames"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret, err := handler.ProfileService.GetProfileIdsByUsernames(ctx, input.Usernames)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(ret)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetPersonalDataByProfileId(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPersonalDataByProfileId-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var id uint
	vars := mux.Vars(r)
	id = util.String2Uint(vars["id"])
	personalData, err := handler.ProfileService.GetPersonalDataByProfileId(ctx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(*personalData)
	if err != nil {
		util.Tracer.LogError(span, err)
		return
	}
}

func (handler *Handler) GetProfileInterests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfileInterests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	vars := mux.Vars(r)
	interests, err := handler.ProfileService.GetProfileInterests(ctx, util.String2Uint(vars["id"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(interests)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func safeRegistrationDto(dto dto.RegistrationDto) dto.RegistrationDto {
	dto.Username = template.HTMLEscapeString(dto.Username)
	dto.Name = template.HTMLEscapeString(dto.Name)
	dto.Surname = template.HTMLEscapeString(dto.Surname)
	dto.Biography = template.HTMLEscapeString(dto.Biography)
	dto.BirthDate = template.HTMLEscapeString(dto.BirthDate)
	dto.Email = template.HTMLEscapeString(dto.Email)
	dto.Gender = template.HTMLEscapeString(dto.Gender)
	dto.Telephone = template.HTMLEscapeString(dto.Telephone)
	dto.WebSite = template.HTMLEscapeString(dto.WebSite)
	interests := dto.InterestedIn
	for i := 0; i < len(dto.InterestedIn); i++ {
		interests[i] = template.HTMLEscapeString(dto.InterestedIn[i])
	}
	return dto
}

func safePersonalDataDto(dto dto.PersonalDataDTO) dto.PersonalDataDTO {
	dto.Username = template.HTMLEscapeString(dto.Username)
	dto.Name = template.HTMLEscapeString(dto.Name)
	dto.Surname = template.HTMLEscapeString(dto.Surname)
	dto.Biography = template.HTMLEscapeString(dto.Biography)
	dto.BirthDate = template.HTMLEscapeString(dto.BirthDate)
	dto.Email = template.HTMLEscapeString(dto.Email)
	dto.Gender = template.HTMLEscapeString(dto.Gender)
	dto.Telephone = template.HTMLEscapeString(dto.Telephone)
	dto.Website = template.HTMLEscapeString(dto.Website)
	return dto
}
