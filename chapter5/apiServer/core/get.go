package core

import (
	"ceph/chapter5/apiServer/heartbeat"
	"ceph/chapter5/apiServer/locate"
	es "ceph/chapter5/lib/ElasticSearch"
	"ceph/chapter5/lib/rs"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//获取对象名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionID := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionID) != 0 {
		version, err = strconv.Atoi(versionID[0])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//从es中获取对象的元数据
	meta, err := es.GetMetadata(object, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//如果哈希值为空则说明被标记为删除
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)

	size, _ := strconv.Atoi(meta.Size)

	//获取hash的读取流
	stream, err := getStream(hash, int64(size))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//恢复分片的写入需要用到temp接口的转正，Close方法用于将写入恢复分片转正
	defer stream.Close()
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

//获取object文件的读取流
func getStream(hash string, size int64) (*rs.RSGetStream, error) {
	//查找哪几台数据节点存了该object的数据分片
	locateInfo, err := locate.Locate(hash)
	if err != nil {
		return nil, err
	}
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail,result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	//如果长度不等于6，说明有部分数据丢失
	//则调用ChooseRandomDataServer随机选取部分用于接收恢复分片的数据服务节点
	if len(locateInfo) != rs.ALL_SHARDS {
		//随机返回多个活跃的数据节点
		dataServers, _ = heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
