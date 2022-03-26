package heartbeat

import (
	"ceph/chapter3/lib/RabbitMQ"
	"os"
	"time"
)

//开始发送心跳包
func StartHeartbeat() {
	//创建一个rabbitmq结构体的实例
	r := RabbitMQ.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer r.Close()
	for {
		//将LISTEN_ADDRESS对应的val,即本地地址发送到apiServers里去
		//所有绑定apiServers这个exchange的RabbitMQ都会收到这个消息
		r.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
