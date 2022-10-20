package objects

import (
	es "OSS/lib/ElasticSearch"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Del(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	defer r.Body.Close()
	//获取object的名称
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	//从es中获取object的最新版本
	version, err := es.SearchLatestVersion(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if version.Version == "" {
		log.Println("Not found", object)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//向es中插入object一个版本，size=0,hash=""，代表着是插入标记
	v, _ := strconv.Atoi(version.Version)
	err = es.PutMetadata(object, v+1, 0, "")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
