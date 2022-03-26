package main

import (
	"ceph/chapter2/apiServer/core"
	"ceph/chapter2/apiServer/heartbeat"
	"ceph/chapter2/apiServer/locate"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	//开始连接apiServers这个exchanges，将数据服务节点的地址保存起来
	go heartbeat.ListenHeartbeat()

	//REST接口 主要是GET和PUT
	//http://apiServerIP/objects/<xxx>
	http.HandleFunc("/objects/", core.Handler)

	//REST接口 主要是找到<xxx>存在于哪个数据节点
	//http://apiServerIP/locate/<xxx>
	http.HandleFunc("/locate/", locate.Handler)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>apiServer Start running <===")

	//监听并启动 ip在tools中规划好了
	//目前是10.29.2.1和10.29.2.2
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
