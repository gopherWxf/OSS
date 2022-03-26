package locate

/*
	http://apiServerIP/locate/<xxx>
	这种形式的url请求过来，api层会向dataServers这个exchange发布object
	如果某个dataServer检测到自己存在这个object，那么它会返回自己的地址给api层
	api层接收到了这个消息，将这个地址消息返回给请求
*/

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	//查找哪台数据节点存了该object
	serverAddr, err := Locate(object)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//将数据节点地址作为相应返回
	bytes, _ := json.Marshal(serverAddr)
	w.Write(bytes)
}
