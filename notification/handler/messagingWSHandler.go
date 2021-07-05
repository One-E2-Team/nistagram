package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"nistagram/notification/dto"
	"nistagram/util"
	"strings"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

var wsMessageMap = make(map[uint][]*websocket.Conn, 0)

func remove(s []*websocket.Conn, conn *websocket.Conn) []*websocket.Conn{
	for key, value := range s {
		if value == conn {
			s[len(s)-1], s[key] = s[key], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func (handler *Handler) MessagingWebSocket(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("MessagingWebSocket-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		util.Tracer.LogError(span, err)
		log.Print("upgrade:", err)
		return
	}
	id := util.GetLoggedUserIDFromToken(request)
	wsMessageMap[id] = append(wsMessageMap[id], conn)
	defer func(c *websocket.Conn) {
		wsMessageMap[id] = remove(wsMessageMap[id], c)
		err := c.Close()
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err.Error())
		}
	}(conn)
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			util.Tracer.LogError(span, err)
			log.Println("read:", err)
			break
		}
		fmt.Printf("%s\n", message)
		var requestData dto.WSRequestBodyDTO
		err = json.Unmarshal(message, &requestData)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			break
		}
		fmt.Println(requestData.Jwt)
		requestData.Jwt = strings.Replace(requestData.Jwt, "\"", "", -1)
		fmt.Println(requestData.Jwt)
		if id != util.GetLoggedUserIDFromPureToken(requestData.Jwt) {
			util.Tracer.LogError(span, fmt.Errorf("possible websocket hijacking, closing"))
			break
		}
		response, ok := handler.messageMultiplexer(ctx, id, requestData.Request, []byte(requestData.Data))
		if ok {
			err = conn.WriteMessage(mt, response)
			if err != nil {
				util.Tracer.LogError(span, err)
				log.Println("write:", err)
				break
			}
		}
	}
}

func (handler *Handler) messageMultiplexer(ctx context.Context, senderId uint, request string, data []byte) ([]byte, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "messageMultiplexer-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", senderId, request))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	switch request {
	case "Random":
		return handler.RndHandler(nextCtx, senderId, data)
	case "SendMessage":
		return handler.SendMessage(nextCtx, senderId, data)
	default:
		return []byte("{\"status\": \"invalid call\"}"), true
	}
}

func TrySendMessage(ctx context.Context, profileId uint, responseType string, data []byte) {
	span := util.Tracer.StartSpanFromContext(ctx, "TrySendMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", profileId, responseType))
	temp := dto.WSResponseBodyDTO{
		Response: responseType,
		Data:     string(data),
	}
	ret, err := json.Marshal(temp)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return
	}
	for _, value := range wsMessageMap[profileId] {
		value.WriteMessage(websocket.TextMessage, ret)
	}
}