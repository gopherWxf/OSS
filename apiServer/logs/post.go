package logs

import (
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Post(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 解析json到map
	body, _ := ioutil.ReadAll(r.Body)
	param := make(map[string]interface{})
	if len(body) != 0 {
		err := json.Unmarshal(body, &param)
		if err != nil {
			golog.Error.Println("json Unmarshal err：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// 从路径中获得分页参数
	index, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	// 通过es获取log
	result := pageHelper(index, getLog(param, w))

	// 序列化json
	b, _ := json.Marshal(result)
	// 写入响应体
	w.Write(b)
}

// 从es轮训获得日志
func getLog(param map[string]interface{}, w http.ResponseWriter) []es.Log {
	// 分页参数
	from := 0
	size := 1000
	result := make([]es.Log, 0)
	//无限循环
	for {
		// 通过对象名 调用es包的SearchAllVersions，返回某个对象的元数据的数组
		searchLog, err := es.SearchLog(param, from, size)
		fmt.Println(len(searchLog), searchLog)
		if from >= 9000 {
			return result
		}
		// 如果报错
		if err != nil {
			// 打印错误并返回500
			golog.Error.Println("es search log err：", err)
			w.WriteHeader(http.StatusInternalServerError)
			return result
		}
		result = append(result, searchLog...)
		// 如果长度数据长度不等于1000，此时没有更多的数据了
		if len(searchLog) != size {
			return result
		}
		//否则把from的值+1000进行下一次迭代
		from += size
	}
}

// 分页结构体
type logInfo struct {
	Size int64    `json:"size"` // 数据总长度
	Data []es.Log `json:"data"` // 每页的数据
}

const SIZE = 20 // 每页显示条数

// 分页
func pageHelper(page int, data []es.Log) logInfo {
	size := len(data)
	info := logInfo{int64(size), nil}
	if size == 0 { //如果长度为0 直接返回
		info.Size = 0
		return info
	}

	log := make([]es.Log, 0)
	start := (page - 1) * SIZE
	end := page * SIZE
	if start > len(data) {
		return info
	}
	if len(data) < end {
		end = len(data)
	}
	log = data[start:end]
	info.Data = log

	return info
}
