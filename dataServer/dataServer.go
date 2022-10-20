package main

import (
	"OSS/dataServer/heartbeat"
	"OSS/dataServer/locate"
	"OSS/dataServer/objects"
	"OSS/dataServer/temp"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func InitRouter(r *gin.Engine) {

	//objects
	{
		r.GET("/objects/:name", objects.Get)
		r.DELETE("/objects/:name", objects.Del)
	}
	//temp
	{
		r.PUT("/temp/:name", temp.Put)
		r.PATCH("/temp/:name", temp.Patch)
		r.POST("/temp/:name", temp.Post)
		r.DELETE("/temp/:name", temp.Del)
		r.HEAD("/temp/:name", temp.Head)
		r.GET("/temp/:name", temp.Get)
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
	//第一次启动时，将所有对象存储到map
	locate.CollectObjects()
	//开始发送心跳包
	go heartbeat.StartHeartbeat()
	//监听来自接口服务local的GET请求,查找本地是否有这个文件,有则发送消息
	go locate.StartLocate()

	r := gin.Default()
	InitRouter(r)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>dataServer Start running <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
