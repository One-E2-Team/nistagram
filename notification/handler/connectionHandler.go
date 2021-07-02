package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/util"
)

func (handler *Handler) GetMessageConnections(w http.ResponseWriter, r *http.Request) {
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	result, err := handler.Service.GetMessageConnections(loggedUserId)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}
