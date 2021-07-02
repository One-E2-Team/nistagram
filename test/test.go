package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"nistagram/notification/dto"
	"os"
)


func main() {
	reader := bufio.NewReader(os.Stdin)
	//jwt, _ := reader.ReadString('\n')
	jwt:= "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRVc2VySWQiOjEsImV4cCI6MTYyNTM0MTEwMywiaWF0IjoxNjI1MjU0NzAzLCJpc3MiOiJhdXRoX3NlcnZpY2UifQ.lTONxXzDgnhnib8ulsf6RIJ4p9alaMzjefXoS9XjGyY"
	header := http.Request{
		Header: map[string][]string{"Authorization": []string{"Bearer " + jwt}},
	}
	c, resp, err := websocket.DefaultDialer.Dial("ws://localhost:8090/messaging", header.Header)
	if resp != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("handshake %s", string(body))
	}
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()


	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s\n", message)
		}
	}()

	for {
		fmt.Println("input request")
		request, _ := reader.ReadString('\n')
		fmt.Println("input data")
		data, _ := reader.ReadString('\n')
		var finalRequest = dto.WSRequestBodyDTO{
			Jwt:     jwt,
			Request: request[:len(request)-1],
			Data:    data[:len(data)-1],
		}
		byteRequest, _ := json.Marshal(finalRequest)
		err := c.WriteMessage(websocket.TextMessage, byteRequest)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}