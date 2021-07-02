package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"nistagram/notification/dto"
	"nistagram/util"
)

var upgrader = websocket.Upgrader{}

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
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	id := util.GetLoggedUserIDFromToken(request)
	wsMessageMap[id] = append(wsMessageMap[id], conn)
	defer func(c *websocket.Conn) {
		wsMessageMap[id] = remove(wsMessageMap[id], c)
		err := c.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(conn)
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var requestData dto.WSRequestBodyDTO
		err = json.Unmarshal(message, &requestData)
		if err != nil {
			fmt.Println(err)
			break
		}
		if id != util.GetLoggedUserIDFromPureToken(requestData.Jwt) {
			break
		}
		response, ok := handler.messageMultiplexer(id, requestData.Request, []byte(requestData.Data))
		if ok {
			err = conn.WriteMessage(mt, response)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func (handler *Handler) messageMultiplexer(senderId uint, request string, data []byte) ([]byte, bool) {
	switch request {
	case "Random":
		return handler.RndHandler(senderId, data)
	default:
		return []byte("{\"status\": \"invalid call\"}"), true
	}
}

func TrySendMessage(profileId uint, data []byte) {
	for _, value := range wsMessageMap[profileId] {
		value.WriteMessage(websocket.TextMessage, data)
	}
}