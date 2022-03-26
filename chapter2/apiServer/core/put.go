package core

import (
	"ceph/chapter2/apiServer/heartbeat"
	"ceph/chapter2/objectStream"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//获取object的名称
	object := strings.Split(r.URL.String(), "/")[2]
	//将body中的内容存入object流中
	statusCode, err := storeObject(r.Body, object)
	if err != nil {
		log.Println(err)
	}
	//返回状态
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
