package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
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
	var req dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while decoding.\"}"))
		return
	}

	v := validator.New()

	_ = v.RegisterValidation("common_pass", func(fl validator.FieldLevel) bool {
		f, err := os.Open("../common_pass.txt")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusOK)
			return false
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if strings.Contains(fl.Field().String(), scanner.Text()) {
				return false
			}
		}
		return true
	})

	_ = v.RegisterValidation("weak_pass", func(fl validator.FieldLevel) bool {
		//^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/])[A-Za-z0-9*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/]{8,}$
		if len(fl.Field().String()) < 8 {
			return false
		}
		ret, _ := regexp.MatchString("(.*[a-z].*)", fl.Field().String())
		if ret == false {
			return false
		}
		ret, _ = regexp.MatchString("(.*[A-Z].*)", fl.Field().String())
		if ret == false {
			return false
		}
		ret, _ = regexp.MatchString("(.*[0-9].*)", fl.Field().String())
		if ret == false {
			return false
		}
		ret, _ = regexp.MatchString("(.*[*!@#$%^&(){}\\[:;\\]<>,.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if err != nil {
			fmt.Println(err)
			return false
		}
		return ret
	})

	_ = v.RegisterValidation("bad_username", func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 3 || len(fl.Field().String()) > 15 {
			return false
		}
		ret, _ := regexp.MatchString("([*!@#$%^&(){}\\[:;\\]<>,.?~+\\-\\\\=|/ ])", fl.Field().String())
		if ret {
			return false
		}
		return true
	})

	err = v.Struct(req)

	if err != nil {
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
	err = handler.ProfileService.Register(req)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while registering.\"}"))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}
	return
}

func (handler *Handler) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := handler.ProfileService.Search(template.HTMLEscapeString(vars["username"]))
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result, err := handler.ProfileService.GetProfileByUsername(template.HTMLEscapeString(vars["username"]))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		return
	}
}

func (handler *Handler) ChangeProfileSettings(w http.ResponseWriter, r *http.Request) {
	var req dto.ProfileSettingsDTO
	methodPath := "nistagram/profile/handler.ChangeProfileSettings"
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	err = handler.ProfileService.ChangeProfileSettings(req, userId)
	if err != nil {
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	util.Logging(util.INFO, methodPath, util.GetIPAddress(r), "Successful profile settings changed, "+util.GetLoggingStringFromID(userId), "profile")
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) ChangePersonalData(w http.ResponseWriter, r *http.Request) {
	var req dto.PersonalDataDTO
	methodPath := "nistagram/profile/handler.ChangePersonalData"
	err := json.NewDecoder(r.Body).Decode(&req)
	req = safePersonalDataDto(req)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "profile")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	oldUsername, oldEmail, err := handler.ProfileService.ChangePersonalData(req, userId)
	if err != nil {
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
	var key string
	err := json.NewDecoder(r.Body).Decode(&key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProfileService.Test(key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetAllInterests(w http.ResponseWriter, _ *http.Request) {
	interests, err := handler.ProfileService.GetAllInterests()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(interests)
	if err != nil {
		return
	}
}

func (handler *Handler) GetMyProfileSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := handler.ProfileService.GetMyProfileSettings(util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		return
	}

}

func (handler *Handler) GetMyPersonalData(w http.ResponseWriter, r *http.Request) {
	data, err := handler.ProfileService.GetMyPersonalData(util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func (handler *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	var id uint
	vars := mux.Vars(r)
	id = util.String2Uint(vars["id"])
	profile, err := handler.ProfileService.GetProfileByID(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(*profile)
	if err != nil {
		return
	}
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
