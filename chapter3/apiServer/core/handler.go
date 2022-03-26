package core

/*
	http://apiServerIP/objects/<xxx>
	这种形式的url相应过来
	如果是PUT：
		首先将r.body里面的内容写入到一个pipe的writer流里面
		http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		将消息转发给上面这个我们定义的请求中
	如果是GET：
		首先去查哪台数据节点存在该object(<xxx>)
		url="http://" + server + "/objects/" + object
		向数据节点发送GET请求，将请求出来的res.body作为一个reader流
		将流中的内容写入到原来的请求里面去作为响应
	如果是DEL:
		//TODO
*/
import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
