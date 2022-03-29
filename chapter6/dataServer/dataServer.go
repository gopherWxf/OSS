package main

import (
	"ceph/chapter6/dataServer/core"
	"ceph/chapter6/dataServer/heartbeat"
	"ceph/chapter6/dataServer/locate"
	"ceph/chapter6/dataServer/temp"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	//第一次启动时，将所有对象存储到map
	locate.CollectObjects()
	//开始发送心跳包
	go heartbeat.StartHeartbeat()
	//监听来自接口服务local的GET请求,查找本地是否有这个文件,有则发送消息
	go locate.StartLocate()

	//REST接口 主要是GET和PUT
	//http://ip/objects/<xxx>
	http.HandleFunc("/objects/", core.Handler)

	//REST接口 主要是PUT,PATCH,POST,DEL
	//http://ip/temp/
	http.HandleFunc("/temp/", temp.Handler)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>dataServer Start running <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
