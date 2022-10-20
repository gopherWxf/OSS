package heartbeat

import (
	RedisMQ "OSS/lib/Redis"
	"os"
	"time"
)

//开始发送心跳包
func StartHeartbeat() {
	rdb := RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer rdb.Client.Close()

	channel := "hearbeat"
	msg := os.Getenv("LISTEN_ADDRESS")

	for {
		//将LISTEN_ADDRESS对应的val,即本地地址发送到apiServers里去
		//所有绑定apiServers这个pubsub的redis都会收到这个消息
		rdb.Publish(channel, msg)
		time.Sleep(time.Second * 5)
	}
}
