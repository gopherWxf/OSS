package temp

import (
	"ceph/chapter7/lib/rs"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取数据节点已经储存该对象多少数据了
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
