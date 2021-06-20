package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/util"
	"strconv"
)

func (handler *Handler) IsBlocked(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok := handler.ConnectionService.IsBlocked(id, uint(id1))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	type resp struct {
		Blocked bool `json:"blocked"`
	}
	json.NewEncoder(w).Encode(resp{Blocked: ok})
}



func (handler *Handler) ToggleBlockProfile(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	_, ok := handler.ConnectionService.ToggleBlock(followerId, profileId)

	if !ok {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("{\"status\":\"ok\"}"))
	writer.Header().Set("Content-Type", "application/json")
}
