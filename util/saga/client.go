package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ChannelHandler struct {
	Channel string
	Handler func(*redis.Client, Message)
}

func SubscribeAndRunPubSubHandlers(client *redis.Client, handlers ...ChannelHandler) {
	if client == nil {
		client = connectToMessageBroker()
	}
	var channels []string
	for _, val := range handlers {
		channels = append(channels, val.Channel)
	}
	pubsub := client.Subscribe(context.TODO(), channels...)

	RunHandlersOnPubSub(client, pubsub, handlers...)

}

func RunHandlersOnPubSub(client *redis.Client, pubsub *redis.PubSub, handlers ...ChannelHandler) {
	if _, err := pubsub.Receive(context.TODO()); err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = pubsub.Close() }()
	ch := pubsub.Channel()

	for{
		select{
		case msg := <-ch:
			m := Message{}
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				fmt.Println(err)
				continue
			}
			for _, val := range handlers {
				if val.Channel == msg.Channel {
					val.Handler(client, m)
				}
			}
		}
	}
}

/*

*/
func SendToReplyChannel(client *redis.Client, m *Message, action string, nextService string, senderService string){
	var err error
	m.Action = action
	m.NextService = nextService
	m.SenderService = senderService
	if err = client.Publish(context.TODO(),ReplyChannel, m).Err(); err != nil {
		fmt.Printf("Error publishing done-message to %s channel", ReplyChannel)
	}
	fmt.Printf("Done message published to channel :%s", ReplyChannel)
}