package RedisMQ

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type PubSubMQ struct {
	Client *redis.Client
}

func NewRedis(redisAddr string) *PubSubMQ {
	client := new(PubSubMQ)
	client.Client = redis.NewClient(&redis.Options{
		Addr: "101.43.17.240:6379",

		//Addr:     redisAddr,
		//Username: "root",
		//Password: "123456",
	})
	return client
}

func (rdb *PubSubMQ) Publish(channel string, message interface{}) {
	err := rdb.Client.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (rdb *PubSubMQ) Subscribe(channels ...string) *redis.PubSub {
	pubsub := rdb.Client.Subscribe(context.Background(), channels...)
	return pubsub
}
