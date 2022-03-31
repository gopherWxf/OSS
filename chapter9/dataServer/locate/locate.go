package locate

/*
	StartLocate会监听dataServers这个exchange
	从中取出api层想要查找的object hash
	在自己目录下查找是否有该文件
	如果有：
		r.Send(msg.ReplyTo, types.LocateMessage{
			Addr: os.Getenv("LISTEN_ADDRESS"),
			Id:   id,
		})
		直接将数据节点本身的地址，发给消息来源的消息队列中
		Send的过程不经过exchange，一对一的
	如果没有：
		不做任何处理
*/
import (
	"ceph/chapter9/lib/RabbitMQ"
	"ceph/chapter9/lib/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
		hash, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		//如果Locate返回！=-1,则说明本地有这个文件,则将这个消息写入dataServers
		id := Locate(hash)
		if id != -1 {
			r.Send(msg.ReplyTo, types.LocateMessage{
				Addr: os.Getenv("LISTEN_ADDRESS"),
				Id:   id,
			})
		}
	}
}

var objects = make(map[string]int)
var mutex sync.Mutex

//查看本地是否有这个文件，并且告知本节点存储的是该对象的哪个分片
func Locate(hash string) int {
	mutex.Lock()
	k, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return k
}
func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}
func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

//第一次启动时，将所有对象存储到map
func CollectObjects() {
	//读取$STORAGE_ROOT/objects/目录里的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		//对这些文件调用filepath.Base获取基本文件名
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		objects[hash] = id
	}
}
