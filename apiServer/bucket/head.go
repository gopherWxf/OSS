package bucket

import (
	es "OSS/lib/ElasticSearch"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 查找bucket
	httpCode := es.SearchBucket(bucket)

	w.WriteHeader(httpCode)
}
