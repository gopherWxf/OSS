package heartbeat

import (
	"ceph/chapter3/lib/RabbitMQ"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//map key是数据节点的地址，val是时间戳
var dataServersMap = make(map[string]time.Time)
var mutex sync.Mutex

//监听apiServers，将数据服务节点的地址保存起来
func ListenHeartbeat() {
	//创建一个rabbitmq结构体的实例
	r := RabbitMQ.NewRabbitMQ(os.Getenv("RABBITMQ_SERVER"))
	defer r.Close()
	//将这个rabbitmq绑定到apiServers这个exchange上
	r.Bind("apiServers")
	//返回一个消费队列消息的channel
	msgCHAN := r.Consume()
	//检查心跳消息，超时就移除数据节点
	go removeExpiredDataServer()
	for msg := range msgCHAN {
		//将数据节点的地址从channal中读取出来
		dataServerAddr, _ := strconv.Unquote(string(msg.Body))
		//存入map中
		mutex.Lock()
		dataServersMap[dataServerAddr] = time.Now()
		mutex.Unlock()
	}
}

//检查心跳消息，超时就移除数据节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for dataServerAddr, timeStamp := range dataServersMap {
			if timeStamp.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServersMap, dataServerAddr)
			}
		}
		mutex.Unlock()
	}
}

//以切片的形式返回所有活跃的数据节点地址
func GetDataServers() (dataServersSlice []string) {
	mutex.Lock()
	for dataServerAddr, _ := range dataServersMap {
		dataServersSlice = append(dataServersSlice, dataServerAddr)
	}
	mutex.Unlock()
	return
}

//随机返回一个活跃的数据节点
func ChooseRandomDataServer() (string, error) {
	dataServersSlice := GetDataServers()
	n := len(dataServersSlice)
	if n == 0 {
		return "", errors.New("not found any dataServer")
	}
	return dataServersSlice[rand.Intn(n)], nil
}
