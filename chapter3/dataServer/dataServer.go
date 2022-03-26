package main

import (
	"ceph/chapter3/dataServer/core"
	"ceph/chapter3/dataServer/heartbeat"
	"ceph/chapter3/dataServer/locate"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	//开始发送心跳包
	go heartbeat.StartHeartbeat()
	//监听来自接口服务local的GET请求,查找本地是否有这个文件,有则发送消息
	go locate.StartLocate()

	//REST接口 主要是GET和PUT
	//http://ip/objects/<xxx>
	http.HandleFunc("/objects/", core.Handler)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>dataServer Start running <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}