package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/util"
)

func (handler *Handler) GetMessageConnections(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMessageConnections-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	result, err := handler.Service.GetMessageConnections(ctx, loggedUserId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}
