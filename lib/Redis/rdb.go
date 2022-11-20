package RedisMQ

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sort"
	"strconv"
)

type RDB struct {
	Client *redis.Client
	//Client *redis.ClusterClient
}

func NewRedis(redisAddr string) *RDB {
	client := new(RDB)
	client.Client = redis.NewClient(&redis.Options{
		//redis访问保护
		Addr: "101.43.17.240:6379",

		//Addr:     redisAddr,
		Username: "root",
		//Password: "123456",
	})
	return client
}

//func NewRedis(redisAddr string) *RDB {
//	client := new(RDB)
//
//	redisClusterAddrsString := os.Getenv("REDIS_CLUSTER")
//	redisClusterAddr := strings.Split(redisClusterAddrsString, ",")
//
//	client.Client = redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs:    redisClusterAddr,
//		Password: os.Getenv("REDIS_PASSWORD"),
//		PoolSize: 20,
//	})
//	_, err := client.Client.Ping(context.Background()).Result()
//	if err != nil {
//		panic(err)
//	}
//	return client
//}

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
func (rdb *RDB) GetEcharts(patten string) (mp map[string]int64) {
	mp = map[string]int64{}
	keys, err := rdb.Client.Keys(context.Background(), patten).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		result, _ := rdb.Client.Get(context.Background(), key).Result()
		i, _ := strconv.Atoi(result)
		mp[key] = int64(i)
	}
	return
}

func (rdb *RDB) GetUpHoldNum(key string) int64 {
	result, err := rdb.Client.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	atoi, _ := strconv.Atoi(result)
	return int64(atoi)
}

func (rdb *RDB) GetOp(hash string, idx int) *Operation {
	ans := &Operation{}
	keys, err := rdb.Client.Keys(context.Background(), hash+"*").Result()
	if err != nil {
		return ans
	}
	sort.Strings(keys)
	idx *= 5
	//确定日期
	if len(keys) > idx {
		keys = keys[idx:]
	}
	if len(keys) > 5 {
		keys = keys[:5]
	}
	//op日期--list->op日期时间       ->op日期时间--string-->op
	for _, key := range keys {
		if len(key) != len("op2022-11-20") {
			continue
		}
		//key == "op2022-11-20"
		data := key[len(hash):]
		//data == "2022-11-20"
		onedata := make([]OpData0, 0)
		opdatetime, _ := rdb.Client.LRange(context.Background(), key, 0, -1).Result()
		for _, datatime := range opdatetime {
			//datatime == op2022-11-2021:53:26
			op, _ := rdb.Client.Get(context.Background(), datatime).Result()
			time := datatime[len("op2022-11-20"):]
			fmt.Println("datatime:", datatime, "time:", time, "op:", op)
			//time == 21:53:26
			onetime := OpData0{
				Operation: op,
				Time:      time,
				Date:      data,
			}
			onedata = append(onedata, onetime)
		}
		alldata := OpData{
			data: data,
			info: onedata,
		}
		ans.OperationData = append(ans.OperationData, alldata)
	}
	ans.OperationSize = int64(len(ans.OperationData))
	return ans
}

type Operation struct {
	OperationSize int64
	OperationData []OpData
}

type OpData struct {
	data string
	info []OpData0
}
type OpData0 struct {
	Operation string
	Time      string
	Date      string
}

func (rdb *RDB) Incr(key string) {
	rdb.Client.Incr(context.Background(), key)
}
func (rdb *RDB) InsertOp(op, date, time string) {
	//op日期--list-->op日期时间       op日期时间--string-->op
	opdata := "op" + date
	opdatatime := opdata + time
	rdb.Client.LPush(context.Background(), opdata, opdatatime)
	rdb.Client.Set(context.Background(), opdatatime, op, 0)
}
