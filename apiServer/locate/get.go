package locate

/*
	http://apiServerIP/locate/<hash>
	这种形式的url请求过来，api层会向dataServers这个exchange发布hash值
	如果某个dataServer检测到自己存在这个hash值，那么它会返回自己的地址给api层
	api层接收到了这个消息，将这个地址消息返回给请求
	注意:1.哈希值可能包含/这个字符，而数据节点存储的都是将/转义后的值 url.PathEscape(hash)
		2.解析url的时候不能再以/分隔，因为hash值里面有可能包含这符号，会导致hash解析错误
*/

import (
	"OSS/lib/golog"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	hash := strings.Split(r.URL.String(), "locate/")[1]
	//查找哪台数据节点存了该hash
	serverAddrs, err := Locate(url.PathEscape(hash))
	if err != nil {
		golog.Error.Println("查找哪台数据节点存了该hash出错 err：", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//将数据节点地址作为相应返回
	bytes, _ := json.Marshal(serverAddrs)
	w.Write(bytes)
}
