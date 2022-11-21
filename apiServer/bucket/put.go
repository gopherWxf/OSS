package bucket

import (
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Put(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	bucket := r.Header.Get("bucket")
	if bucket == "" {
		golog.Error.Println("请求头缺少bucket字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := es.AddBucket(bucket)
	if err != nil {
		golog.Error.Println("增加bucket出错：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	golog.Info.Println("增加bucket：", bucket, " 成功")
	w.WriteHeader(http.StatusCreated)
}
