package objects

import (
	"ceph/chapter9/apiServer/heartbeat"
	"ceph/chapter9/apiServer/locate"
	es "ceph/chapter9/lib/ElasticSearch"
	"ceph/chapter9/lib/rs"
	"ceph/chapter9/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//获取对象名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	//获取对象内容的大小
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取hash值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//如果已经存在，则直接添加元数据到ES中，并返回200ok
	if locate.Exist(url.PathEscape(hash)) {
		err = es.AddVersion(object, hash, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	//如果不存在，则随机挑选6个数据节点
	dataServers, err := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	//调用rs.NewRSResumablePutStream生成数据流stream
	stream, err := rs.NewRSResumablePutStream(dataServers, object, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//调用ToToken生成一个字符串token，放入Location响应头中，返回201StatusCreated
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	//TODO base64包含/ 而pathEscape会转化/ ,那么解析的时候怎么办呢，用	url.PathUnescape()
	w.WriteHeader(http.StatusCreated)
}
