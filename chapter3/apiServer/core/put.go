package core

import (
	"ceph/chapter3/apiServer/heartbeat"
	es "ceph/chapter3/lib/ElasticSearch"
	"ceph/chapter3/lib/objectStream"
	"ceph/chapter3/utils"
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
	//将hash值作为数据节点存储文件的名称，实现对象名与内容的解耦
	statusCode, err := storeObject(r.Body, url.PathEscape(hash))
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
	//获取头部内容的长度
	size := utils.GetSizeFromHeader(r.Header)
	//给该对象增加新版本
	err = es.AddVersion(object, hash, size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(statusCode)
}

//将body中的内容存入object流中
func storeObject(r io.Reader, object string) (int, error) {
	//随机选择一个数据节点创建一个输入流
	stream, err := putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	//将put的数据写入到输入流中
	io.Copy(stream, r)
	//关闭的时候会读一个错误日志，如果有错误说明数据转发给数据节点的时候有错误发送
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

//随机选择一个数据节点创建一个输入流
func putStream(object string) (*objectStream.PutStream, error) {
	//随机选一个活跃节点
	serverAddr, err := heartbeat.ChooseRandomDataServer()
	if err != nil {
		return nil, err
	}
	//创建一个输入流并返回
	return objectStream.NewPutStream(serverAddr, object), nil
}
