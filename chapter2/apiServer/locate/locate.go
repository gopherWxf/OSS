package locate

import (
	"ceph/chapter2/RabbitMQ"
	"errors"
	"os"
	"strconv"
	"time"
)

//查找哪台数据节点存了该object
func Locate(object string) (string, error) {
	//创建一个rabbitmq结构体的实例
	r := RabbitMQ.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer r.Close()
	//将object发布到dataServers这个exchange里面去，供别人接收
	//所有绑定dataServers这个exchange的队列都会收到这个消息
	r.Publish("dataServers", object)
	//获取消费队列的channel
	ch := r.Consume()
	//等待一秒钟
	go func() {
		time.Sleep(1 * time.Second)
		r.Close()
	}()
	msg := <-ch
	if len(string(msg.Body)) == 0 {
		return "", errors.New("All the data servers were not found object:" + object)
	}
	s, _ := strconv.Unquote(string(msg.Body))
	return s, nil
}
