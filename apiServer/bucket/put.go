package bucket

import (
	es "OSS/lib/ElasticSearch"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Put(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	bucket := r.Header.Get("bucket")
	if bucket == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := es.AddBucket(bucket)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
