package saga

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

func connectToMessageBroker() *redis.Client {
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
			Password: "helloworld",
			DB:       0,
		})

		if err := client.Ping(context.TODO()).Err(); err != nil {
			fmt.Println("Cannot connect to redis! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to redis.")
			break
		}
	}
	return client
}