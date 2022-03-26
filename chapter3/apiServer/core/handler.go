package core

/*
	http://apiServerIP/objects/<xxx>
	这种形式的url相应过来
	如果是PUT：
		获取请求头部中内容的的hash值
		hash := utils.GetHashFromHeader(r.Header)
		将hash值作为数据节点存储文件的名称，实现对象名与内容的解耦
		statusCode, err := storeObject(r.Body, url.PathEscape(hash))
		给该对象增加新版本
		es.AddVersion(object, hash, size)
		PUT完成,返回响应
		w.WriteHeader(statusCode)
	如果是GET：
		从es中获取对象的元数据
		meta, err := es.GetMetadata(object, version)
		如果哈希值为空则说明被标记为删除
		获取hash的读取流
		stream, err := getStream(hash)
		将流中的内容写入到Response中，作为相应内容
		io.Copy(w, stream)
	如果是DEL:
		首先先去获取对象的最新版本信息
		version, err := es.SearchLatestVersion(object)
		es中插入一个新版本版本，size=0,hash=""，标志着删除
		es.PutMetadata(object, v+1, 0, "")
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
