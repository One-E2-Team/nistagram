package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

const (
	ProfileChannel string = "ProfileChannel"
	AuthChannel string = "AuthChannel"
	ConnectionChannel string = "ConnectionChannel"
	PostChannel string = "PostChannel"
	ReplyChannel    string = "ReplyChannel"
	ActionStart     string = "Start"
	ActionRollback  string = "RollbackMsg"
	ActionDone      string = "DoneMsg"
	ActionError     string = "ErrorMsg"
	ChangeProfilesPrivacy string = "ChangeProfilesPrivacy"
	RegisterProfile string = "RegisterProfile"
	DeleteProfile string = "DeleteProfile"
	ProfileService string = "ProfileService"
	AuthService string = "AuthService"
	ConnectionService string = "ConnectionService"
	PostService string = "PostService"
)

type Orchestrator struct{
	Client *redis.Client
	PubSub *redis.PubSub
}

func NewOrchestrator() *Orchestrator{
	var client *redis.Client
	time.Sleep(5 * time.Second)
	var redisHost, redisPort = "localhost", "6379"          // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")        // production environment
	if ok {
		redisHost = "message_broker"
		redisPort = "6379"
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			redisHost = "message_broker"
			redisPort = "6379"
		}
	}
	for {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: "",
			DB:       0,
		})

		if err := client.Ping(context.TODO()).Err(); err != nil {
			fmt.Println("Cannot connect to redis! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Orchestrator connected to redis.")
			break
		}
	}

	orch := &Orchestrator{
		Client: client,
		PubSub: client.Subscribe(context.TODO(), ProfileChannel, AuthChannel,
			ConnectionChannel, PostChannel, ReplyChannel),
	}

	return orch
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
				if m.Action == ActionError{
					o.Rollback(m)
					continue
				}
				if m.Action == ActionDone{
					fmt.Println("Functionality done.")
				}
			}
		}

	}
}

func (o Orchestrator) Next(channel, nextService string, message Message) {
	var err error
	message.Action = ActionStart
	message.NextService = nextService
	if err = o.Client.Publish(context.TODO(),channel, message).Err(); err != nil {
		fmt.Printf("Error publishing start-message to %s channel", channel)
	}
	fmt.Printf("Start message published to channel :%s", channel)
}

func (o Orchestrator) Rollback(message Message) {
	var err error
	var channel string
	switch message.NextService {
	case ProfileService:
		channel = ProfileChannel
	}
	message.Action = ActionRollback
	if err = o.Client.Publish(context.TODO(),channel, message).Err(); err != nil {
		fmt.Printf("Error publishing rollback message to %s channel", ProfileChannel)
	}
}