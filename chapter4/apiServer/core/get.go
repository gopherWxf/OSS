package core

import (
	"ceph/chapter4/apiServer/locate"
	es "ceph/chapter4/lib/ElasticSearch"
	"ceph/chapter4/lib/objectStream"
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
	//获取hash的读取流
	stream, err := getStream(hash)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer stream.(*objectStream.GetStream).Reader.(io.ReadCloser).Close()
	//将流中的内容写入到Response中，作为相应内容
	io.Copy(w, stream)
}

//获取object文件的读取流
func getStream(object string) (io.Reader, error) {
	//查找哪台数据节点存了该object
	serverAddr, err := locate.Locate(object)
	if err != nil {
		return nil, err
	}
	return objectStream.NewGetStream(serverAddr, object)
}
