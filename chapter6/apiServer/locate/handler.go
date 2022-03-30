package locate

/*
	http://apiServerIP/locate/<hash>
	这种形式的url请求过来，api层会向dataServers这个exchange发布hash值
	如果某个dataServer检测到自己存在这个hash值，那么它会返回自己的地址给api层
	api层接收到了这个消息，将这个地址消息返回给请求
	注意:1.哈希值可能包含/这个字符，而数据节点存储的都是将/转义后的值 url.PathEscape(hash)
		2.解析url的时候不能再以/分隔，因为hash值里面有可能包含这符号，会导致hash解析错误

	注意：该接口已无意义，因为做了数据的冗余备份，那么一个对象文件会被分成6个数据文件分片存储在不同的数据节点上

*/

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	hash := strings.Split(r.URL.String(), "locate/")[1]
	//查找哪台数据节点存了该hash
	serverAddrs, err := Locate(url.PathEscape(hash))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//将数据节点地址作为相应返回
	bytes, _ := json.Marshal(serverAddrs)
	w.Write(bytes)
}
