package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)



type Orchestrator struct{
	Client *redis.Client
	PubSub *redis.PubSub
}

func NewOrchestrator(channels []string) *Orchestrator {
	client := connectToMessageBroker()
	orch := &Orchestrator{
		Client: client,
		PubSub: client.Subscribe(context.TODO(), channels...),
	}
	return orch
}

func (o Orchestrator) Next(channel, nextService string, message Message) {
	var err error
	message.Action = ActionStart
	message.NextService = nextService
	if err = o.Client.Publish(context.TODO(), channel, message).Err(); err != nil {
		fmt.Printf("Error publishing start-message to %s channel", channel)
	}
	fmt.Printf("Start message published to channel :%s", channel)
}

func (o Orchestrator) Rollback(m Message) {
	var err error
	var channel string
	switch m.NextService {
	case ProfileService:
		channel = ProfileChannel
	}
	m.Action = ActionRollback
	if err = o.Client.Publish(context.TODO(), channel, m).Err(); err != nil {
		fmt.Printf("Error publishing rollback message to %s channel", ProfileChannel)
	}
}

func (o Orchestrator) Start(){
	var err error
	if _, err = o.PubSub.Receive(context.TODO()); err != nil{
		fmt.Println(err)
		return
	}

	ch := o.PubSub.Channel()
	defer o.PubSub.Close()

	fmt.Println("Starting orchestrator..")

	for{
		select{
		case msg := <-ch:
			m := Message{}
			if err = json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				fmt.Println(err)
				continue
			}

			switch msg.Channel {
			case ReplyChannel:
				if m.Action != ActionDone{
					o.Rollback(m)
					continue
				}else{
					fmt.Println("Functionality done.")
				}
			}
		}

	}
}