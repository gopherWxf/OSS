package locate

/*
	StartLocate会监听dataServers这个exchange
	从中取出api层想要查找的object
	在自己目录下查找是否有该文件
	如果有：
		r.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		直接将数据节点本身的地址，发给消息来源的消息队列中
		Send的过程不经过exchange，一对一的
	如果没有：
		不做任何处理
*/
import (
	"ceph/chapter2/RabbitMQ"
	"os"
	"strconv"
)

//监听来自接口服务local的请求,查找本地是否有这个文件,有则发送消息
func StartLocate() {
	//创建一个rabbitmq结构体的实例
	r := RabbitMQ.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer r.Close()
	//绑定dataServers这个exchange
	r.Bind("dataServers")
	//获取dataServers发来的消息的消息队列channel
	ch := r.Consume()
	//遍历消息队列
	for msg := range ch {
		//解析出object名字
		object, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		//如果Locate返回true,则说明本地有这个文件,则将这个消息写入dataServers
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			r.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

//查看本地是否有这个文件
func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
