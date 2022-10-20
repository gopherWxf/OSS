package heartbeat

import (
	RedisMQ "OSS/lib/Redis"
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

//map key是数据节点的地址，val是时间戳
var dataServersMap = make(map[string]time.Time)
var mutex sync.Mutex

//监听apiServers，将数据服务节点的地址保存起来
func ListenHeartbeat() {
	//创建一个redis连接
	rdb := RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer rdb.Client.Close()

	channel := "heartbeat"
	pubsub := rdb.Subscribe(channel)

	//检查心跳消息，超时就移除数据节点
	go removeExpiredDataServer()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Println("redis err:", err)
		}
		dataServerAddr := msg.Payload
		mutex.Lock()
		dataServersMap[dataServerAddr] = time.Now()
		mutex.Unlock()
		//log.Println("recv data addr",dataServerAddr)
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
	for dataServerAddr := range dataServersMap {
		dataServersSlice = append(dataServersSlice, dataServerAddr)
	}
	mutex.Unlock()
	return
}

//随机返回多个活跃的数据节点
//n表示我们需要多少个随机数据节点
//exclude作用是要求返回的节点中不能包含map中的节点，因为数据修复的时候需要排除掉已有的分片所在的节点
func ChooseRandomDataServer(n int, exclude map[int]string) (ds []string, err error) {
	//候选数组
	candidates := make([]string, 0)
	//翻转excluemap，这样排除的时候方便
	reverseExcludeMap := make(map[string]int)
	for k, v := range exclude {
		reverseExcludeMap[v] = k
	}
	//获得所有活跃的数据节点
	dataServersSlice := GetDataServers()
	for _, v := range dataServersSlice {
		_, exclude := reverseExcludeMap[v]
		//这个该节点不再排除map中，则加入候选数组
		if !exclude {
			candidates = append(candidates, v)
		}
	}
	length := len(candidates)
	//如果没有n个数据服务节点，则返回err
	if length < n {
		return nil, errors.New("can not find enough dataServer")
	}
	//打乱，乱序
	p := rand.Perm(length)
	//取前n个数据节点
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return
}
