package temp

import (
	"OSS/lib/golog"
	"OSS/lib/rs"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Head(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	token := strings.Split(r.URL.EscapedPath(), "/")[3]
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		golog.Error.Println("new rs put stream err：", err)
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
	golog.Info.Println("获取数据节点已经储存该对象多少数据成功")
}
