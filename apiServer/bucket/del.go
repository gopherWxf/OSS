package bucket

import (
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Del(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 获得桶名
	bucket := r.Header.Get("bucket")
	if bucket == "" {
		golog.Error.Println("请求头缺少 bucket 字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := es.DelBucket(bucket)
	if err != nil {
		golog.Error.Println("删除 bucket 时出错：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	golog.Info.Println("删除 bucket ：成功", bucket)
	w.WriteHeader(http.StatusOK)
}
