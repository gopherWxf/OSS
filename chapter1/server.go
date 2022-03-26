package main

import (
	"ceph/chapter1/core"
	"net/http"
	"os"
)

func main() {
	//回调函数core.Handler
	http.HandleFunc("/objects/", core.Handler)
	//os.Getenv()检索由键命名的环境变量的值
	http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil)
}
