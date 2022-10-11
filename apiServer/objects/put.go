package objects

import (
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/rs"
	"OSS/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//获取请求头部中内容的的hash值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//获取内容的大小
	size := utils.GetSizeFromHeader(r.Header)
	//将hash值作为数据节点存储文件的名称，实现对象名与内容的解耦
	statusCode, err := storeObject(r.Body, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(statusCode)
		return
	}
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}
	//获取object的名称
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	//给该对象增加新版本
	err = es.AddVersion(object, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(statusCode)
}

//将body中的内容存入object流中  在客户端将内容传递给stream的同时，进行数据校验，如果不一致，则不转正，调用DELETE
func storeObject(r io.Reader, hash string, size int64) (int, error) {
	//是否有节点存过该hash,如果存在，则跳过后续上传，直接返回200ok
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	//随机选择多个数据节点创建一个输入流
	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	//类似于linux的tee命令，参数是io.Reader和io.Writer，返回io.Reader
	//当reader被读取时，实际内容来自于r，被读取中，同时会写入stream
	//即当utils.CalculateHash(reader)从reader中读取数据的同时也写入了stream
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	//如果计算的hash值与客户端传递的hash不一致，则删除临时对象
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	//如果一致，则将临时对象转正
	stream.Commit(true)
	return http.StatusOK, nil
}

//随机选择一个数据节点创建一个输入流
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	//随机选多个活跃节点
	serverAddrs, err := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if err != nil {
		return nil, err
	}
	//创建一个输入流并返回
	return rs.NewRSPutStream(serverAddrs, hash, size)
}
