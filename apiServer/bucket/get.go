package bucket

import (
	"OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// 分页结构体
type bucketInfo struct {
	Size int64    `json:"size"` // 数据总长度
	Data []string `json:"data"` // 每页的数据
}

const SIZE = 4 // 每页显示条数
func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 从路径中获得分页参数
	index, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	// 获得桶名
	bucket := r.Header.Get("bucket")

	buckets := es.GetAllBucket()

	if bucket != "" { // 桶名不为空 则是查询单个
		golog.Info.Println("查询bucket：", bucket)
		var result = make([]string, 0)
		for _, curBucket := range buckets {
			if strings.Contains(curBucket, bucket) {
				result = append(result, curBucket)
			}
		}

		helper := pageHelper(index, result) // 分页
		marshal, _ := json.Marshal(helper)
		w.WriteHeader(http.StatusOK)
		w.Write(marshal)
		return
	}

	// 否则是查询全部
	golog.Info.Println("查询所有 bucket")
	helper := pageHelper(index, buckets) // 分页
	marshal, _ := json.Marshal(helper)

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}

// 比较前一个字符串是否与后一个相同
func strVagueQuery(a string, b string) bool {
	return strings.Contains(b, a)
}

// 分页
func pageHelper(page int, data []string) bucketInfo {
	size := len(data)
	info := bucketInfo{int64(size), nil}
	if size == 0 { //如果长度为0 直接返回
		info.Size = 0
		return info
	}

	metadata := make([]string, 0)
	start := (page - 1) * SIZE
	end := page * SIZE
	if start > len(data) {
		fmt.Println([]int{})
		return info
	}
	if len(data) < end {
		end = len(data)
	}

	metadata = data[start:end]
	info.Data = metadata

	return info
}
