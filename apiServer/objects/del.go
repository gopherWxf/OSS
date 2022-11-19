package objects

import (
	es "OSS/lib/ElasticSearch"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func Del(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucket == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 获取对象名称
	name := strings.Split(r.URL.EscapedPath(), "/")[3]

	//从es中获取object的最新版本
	version, err := es.SearchLatestVersion(bucket, name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if version.Version == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//向es中插入object一个版本，size=0,hash=""，代表着是插入标记
	err = es.PutMetadata(bucket, name, version.Version+1, 0, "")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
