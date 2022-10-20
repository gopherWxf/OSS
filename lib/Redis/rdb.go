package RedisMQ

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RDB struct {
	Client *redis.Client
}

func NewRedis(redisAddr string) *RDB {
	client := new(RDB)
	client.Client = redis.NewClient(&redis.Options{
		Addr: "101.43.17.240:6379",

		//Addr:     redisAddr,
		//Username: "root",
		//Password: "123456",
	})
	return client
}

func (rdb *RDB) Publish(channel string, message interface{}) {
	err := rdb.Client.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (rdb *RDB) Subscribe(channels ...string) *redis.PubSub {
	pubsub := rdb.Client.Subscribe(context.Background(), channels...)
	return pubsub
}
func (rdb *RDB) RemoveFile(hash, ip string) {
	rdb.Client.ZRem(context.Background(), hash, ip)
}
func (rdb *RDB) GetZsetIdAndIP(hash string) ([]string, error) {
	return rdb.Client.ZRange(context.Background(), hash, 0, -1).Result()
}
