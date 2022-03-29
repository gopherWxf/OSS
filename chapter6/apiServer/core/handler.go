package core

/*
	http://apiServerIP/objects/<xxx>
	这种形式的url相应过来
	如果是PUT：
		如果有节点存过该hash,如果存在，则跳过后续上传，直接往es中添加元数据即可
		如果没有，随机选择一个数据节点创建一个输入流
		stream, err := putStream(url.PathEscape(hash), size)
		类似于linux的tee命令，参数是io.Reader和io.Writer，返回io.Reader
		当reader被读取时，实际内容来自于r，被读取中，同时会写入stream
		即当utils.CalculateHash(reader)从reader中读取数据的同时也写入了stream,会触发stream的Write方法
			在新创建steam对象的时候 http://server/temp/hash  POST 见dataServer POST
			request, err := http.NewRequest("POST", "http://"+server+"/temp/"+hash, nil)
			在触发Write方法的时候 http://Server/temp/+w.Uuid PATCH 见dataServer PATCH
		reader := io.TeeReader(r, stream)
		d := utils.CalculateHash(reader)
		如果计算的hash值与客户端传递的hash不一致，则删除临时对象，一致则转正临时对象
		stream.Commit(true) or stream.Commit(false)
	如果是GET：
		首先从es中获取对象的元数据
		meta, err := es.GetMetadata(object, version)
		获取hash的读取流
		stream, err := getStream(hash)
		写入响应中
	如果是DEL:
		从es中获取object的最新版本
		version, err := es.SearchLatestVersion(object)
		向es中插入object一个版本，size=0,hash=""，代表着是插入标记
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
	//TODO 新增post，待写描述注释
	if m == http.MethodPost {
		post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
