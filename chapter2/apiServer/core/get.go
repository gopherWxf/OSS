package core

import (
	"ceph/chapter2/apiServer/locate"
	"ceph/chapter2/objectStream"
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	object := strings.Split(r.URL.String(), "/")[2]
	//获取object文件的读取流
	stream, err := getStream(object)
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
