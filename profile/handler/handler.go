package handler

import (
	"bufio"
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
	var req dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while decoding.\"}"))
		return
	}

	v := validator.New()

	handler.checkCommonPass(w, v)

	handler.checkWeakPass(v, err)

	handler.checkUsername(v)

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
}

func (handler *Handler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	var req dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while decoding.\"}"))
		return
	}

	v := validator.New()

	handler.checkCommonPass(w, v)

	handler.checkWeakPass(v, err)

	handler.checkUsername(v)

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
	err = handler.ProfileService.RegisterAgent(req)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"Server error while registering.\"}"))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}
}

func (handler *Handler) checkUsername(v *validator.Validate) {
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
}

func (handler *Handler) checkWeakPass(v *validator.Validate, err error) {
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
}

func (handler *Handler) checkCommonPass(w http.ResponseWriter, v *validator.Validate) {
	_ = v.RegisterValidation("common_pass", func(fl validator.FieldLevel) bool {
		f, err := os.Open("../common_pass.txt")
		if err != nil {
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
				return false
			}
		}
		return true
	})
}

func (handler *Handler) CreateVerificationRequest(w http.ResponseWriter, r *http.Request) {
	profileId := util.GetLoggedUserIDFromToken(r)
	err := r.ParseMultipartForm(0)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	var reqDto dto.VerificationRequestDTO
	data := r.MultipartForm.Value["data"]
	err = json.Unmarshal([]byte(data[0]), &reqDto)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	picture, picHeader, err := r.FormFile("picture")
	uid := uuid.NewString()

	fileSplitted := strings.Split(picHeader.Filename, ".")
	fileName := uid + "." + fileSplitted[1]

	err = handler.ProfileService.CreateVerificationRequest(profileId, reqDto, fileName)

	if err != nil {
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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	_, err = io.Copy(f, picture)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
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

func (handler *Handler) SearchForTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loggedUserId := util.GetLoggedUserIDFromToken(r)
	result, err := handler.ProfileService.SearchForTag(loggedUserId, template.HTMLEscapeString(vars["username"]))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
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

func (handler *Handler) UpdateVerificationRequest(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyDTO
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = handler.ProfileService.UpdateVerificationRequest(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
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

func (handler *Handler) GetAllCategories(w http.ResponseWriter, _ *http.Request) {
	categories, err := handler.ProfileService.GetAllCategories()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(categories)
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

func (handler *Handler) GetVerificationRequests(w http.ResponseWriter, _ *http.Request) {
	requests, err := handler.ProfileService.GetVerificationRequests()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		return
	}
	fmt.Println(requests, err)
}

func (handler *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	profileId := util.String2Uint(vars["id"])
	err := handler.ProfileService.DeleteProfile(profileId)

	if err != nil {
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
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	err := handler.ProfileService.SendAgentRequest(loggedUserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetAgentRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := handler.ProfileService.GetAgentRequests()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(requests)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetProfileUsernamesByIDs(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Ids []string `json:"ids"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret, err := handler.ProfileService.GetProfileUsernamesByIDs(input.Ids)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(ret)
	if err != nil {
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
	type data struct {
		Interests []string `json:"interests"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profiles, err := handler.ProfileService.GetByInterests(input.Interests)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(profiles)
	if err != nil {
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
	var input dto.ProcessAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := handler.ProfileService.ProcessAgentRequest(input)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetProfileIdsByUsernames(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Usernames []string `json:"usernames"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret, err := handler.ProfileService.GetProfileIdsByUsernames(input.Usernames)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(ret)
	if err != nil {
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
	var id uint
	vars := mux.Vars(r)
	id = util.String2Uint(vars["id"])
	personalData, err := handler.ProfileService.GetPersonalDataByProfileId(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(*personalData)
	if err != nil {
		return
	}
}

func (handler *Handler) GetProfileInterests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	interests, err := handler.ProfileService.GetProfileInterests(util.String2Uint(vars["id"]))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(interests)
	if err != nil {
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
