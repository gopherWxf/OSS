package main

import (
	"OSS/chapter10/apiServer/heartbeat"
	"OSS/chapter10/apiServer/locate"
	"OSS/chapter10/apiServer/objects"
	"OSS/chapter10/apiServer/temp"
	"OSS/chapter10/apiServer/versions"
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
	//http://apiServerIP/objects/<xxx>  这里的<xxx>是对象名
	//http://apiServerIP/objects/<xxx>？version=n
	http.HandleFunc("/objects/", objects.Handler)

	//REST接口 主要是找到<xxx>存在于哪个数据节点
	//http://apiServerIP/locate/<xxx>   这里的<xxx>是hash值
	http.HandleFunc("/locate/", locate.Handler)

	//REST接口 主要是找到<xxx>存在于哪个数据节点
	//http://apiServerIP/versions/ 查看所有对象的所有版本信息
	//http://apiServerIP/versions/<xxx> 查看指定对象的所有版本信息
	http.HandleFunc("/versions/", versions.Handler)

	//REST接口 TODO
	//http://apiServerIP/temp/<xxx>
	http.HandleFunc("/temp/", temp.Handler)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>apiServer Start running <===")

	//监听并启动 ip在tools中规划好了
	//目前是10.29.2.1和10.29.2.2
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
