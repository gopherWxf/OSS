package tools

import (
	es "OSS/lib/ElasticSearch"
	"OSS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	删除过期元数据的工具
	同一对象，最多留存5个历史版本，最早的版本会被删掉
*/

func DelOldMetaDate(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	// 保留几个版本
	version, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	buckets := es.GetAllBucket()
	from, size := 0, 1000
	mpMax := map[string]int{}
	mpMin := map[string]int{}
	for _, bucket := range buckets {
		metas, err := es.SearchAllVersions(bucket, "", from, size)
		if err != nil {
			return
		}
		for _, k := range metas {
			if _, ok := mpMax[k.Name]; ok {
				mpMax[k.Name] = max(mpMax[k.Name], k.Version)
			} else {
				mpMax[k.Name] = k.Version
			}
			if _, ok := mpMin[k.Name]; ok {
				mpMin[k.Name] = min(mpMin[k.Name], k.Version)
			} else {
				mpMin[k.Name] = k.Version
			}
		}
		for k := range mpMax {
			cur := mpMax[k] - mpMin[k] + 1
			if cur > version {
				for v := mpMax[k] - version; v >= mpMin[k]; v-- {
					es.DelMetadata(bucket, k, v)
				}
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	rdb := utils.Rds
	rdb.Incr("OssUpHold")
	op := fmt.Sprintf("进行了保留对象版本操作：保留了全部对象的%d个版本", version)
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("15:04:05")
	rdb.InsertOp(op, date, time)

}
func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
