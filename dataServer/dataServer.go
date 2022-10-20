package main

import (
	"OSS/dataServer/heartbeat"
	"OSS/dataServer/locate"
	"OSS/dataServer/objects"
	"OSS/dataServer/temp"
	RedisMQ "OSS/lib/Redis"
	utils2 "OSS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func InitRouter(r *gin.Engine) {
	//objects
	{
		r.GET("/objects/*id", objects.Get)
		r.DELETE("/objects/*id", objects.Del)
	}
	//temp
	{
		r.PUT("/temp/*id", temp.Put)
		r.PATCH("/temp/*id", temp.Patch)
		r.POST("/temp/*id", temp.Post)
		r.DELETE("/temp/*id", temp.Del)
		r.HEAD("/temp/*id", temp.Head)
		r.GET("/temp/*id", temp.Get)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	utils2.Rds = RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer utils2.Rds.Client.Close()

	//第一次启动时，将所有对象存储到map
	locate.CollectObjects()
	//开始发送心跳包
	go heartbeat.StartHeartbeat()

	////监听来自接口服务local的GET请求,查找本地是否有这个文件,有则发送消息
	//go locate.StartLocate()

	r := gin.Default()
	InitRouter(r)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>dataServer Start running <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	log.Fatal(r.Run(os.Getenv("LISTEN_ADDRESS")))
}
