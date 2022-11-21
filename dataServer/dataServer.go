package main

import (
	"OSS/dataServer/heartbeat"
	"OSS/dataServer/locate"
	"OSS/dataServer/objects"
	"OSS/dataServer/system"
	"OSS/dataServer/temp"
	RedisMQ "OSS/lib/Redis"
	"OSS/lib/golog"
	utils2 "OSS/utils"
	"github.com/gin-gonic/gin"
	"os"
	"time"
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
	//system
	{
		r.GET("/systemInfo", system.Get)
	}
}

func main() {
	// 实时读取日志
	go golog.ReadLog(time.Now().Format("2006-01-02"))

	utils2.Rds = RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer utils2.Rds.Client.Close()

	//第一次启动时，将所有对象存储到map
	locate.CollectObjects()
	//开始发送心跳包
	go heartbeat.StartHeartbeat()

	r := gin.Default()
	InitRouter(r)

	golog.Info.Println(os.Getenv("LISTEN_ADDRESS"), "===>dataServer Start running <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	golog.Info.Println(r.Run(os.Getenv("LISTEN_ADDRESS")))
}
