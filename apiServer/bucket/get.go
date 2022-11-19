package bucket

import (
	"OSS/lib/ElasticSearch"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 从路径中获得分页参数
	index, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	// 获得桶名
	bucket := r.Header.Get("bucket")

	mapping := es.GetAllMapping()

	if bucket != "" { // 桶名不为空 则是查询单个
		myLog.Info.Println(fmt.Sprintf("查询桶 %s", bucket))
		unescape, _ := url.QueryUnescape(bucket)
		var result = make([]string, 0)
		for _, m := range mapping {
			if strVagueQuery(unescape, m) {
				result = append(result, m)
			}
		}

		helper := pageHelper(index, result) // 分页

		marshal, _ := json.Marshal(helper)
		w.WriteHeader(http.StatusOK)
		w.Write(marshal)
		return
	}

	// 否则是查询全部

	helper := pageHelper(index, mapping) // 分页
	marshal, _ := json.Marshal(helper)

	myLog.Info.Println(fmt.Sprintf("查询全部桶，第 %d 页", index))
	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
}
