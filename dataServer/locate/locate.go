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
	"OSS/lib/golog"
	utils2 "OSS/utils"
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//第一次启动时，将所有对象存储到map
func CollectObjects() {
	//读取$STORAGE_ROOT/objects/目录里的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	//写入redis
	rdb := utils2.Rds
	for i := range files {
		//对这些文件调用filepath.Base获取基本文件名
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			golog.Error.Println("read rile err")
			panic(files[i])
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		//objects[hash] = id

		// ZAdd Redis `ZADD key score member [score member ...]` command.
		err = rdb.Client.ZAdd(context.Background(), hash, &redis.Z{
			Score:  float64(id),
			Member: os.Getenv("LISTEN_ADDRESS"),
		}).Err()
		if err != nil {
			golog.Error.Println("redis zadd err：", err)
		}
		golog.Info.Println(hash, float64(id), os.Getenv("LISTEN_ADDRESS"))
	}
}
