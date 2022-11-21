package bucket

import (
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Head(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 获得桶名
	bucket := r.Header.Get("bucket")

	if bucket == "" {
		golog.Error.Println("请求头缺少bucket字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 查找bucket
	golog.Info.Println("查询bucket：", bucket, ",是否存在")
	httpCode := es.SearchBucket(bucket)

	w.WriteHeader(httpCode)
}
