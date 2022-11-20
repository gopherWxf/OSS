package versions

/*
	http://apiServerIP/versions/ 查看所有对象的所有版本信息
	http://apiServerIP/versions/<xxx> 查看指定对象的所有版本信息
	通过es的api去构造url，返回的是json，解析到结构体，然后输出即可
	metas, err := es.SearchAllVersions(object, from, size)
	url := fmt.Sprintf("http://%s/metadata/objects/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
*/
import (
	es "OSS/lib/ElasticSearch"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer ctx.Request.Body.Close()
	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 获取对象名称
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	//%E6%B5%8B%E8%AF%952
	if bucket == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	from, size := 0, 1000
	fmt.Println(name)
	for {
		//获取元数据信息
		metas, err := es.SearchAllVersions(bucket, name, from, size)
		if err != nil {
			log.Println(err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		bytes, _ := json.Marshal(metas)
		w.Write(bytes)
		//如果长度不等于size，说明没有更多的数据了
		if len(metas) != size {
			return
		}
		from += size
	}
}
func AllGet(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer ctx.Request.Body.Close()
	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucket == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 从路径中获得对象名
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	// 获得page
	pageIndex := r.URL.Query()["page"]
	page := 1
	if len(pageIndex) != 0 {
		page, _ = strconv.Atoi(pageIndex[0])
	}
	// 如果有参数 调用es包的SearchApiVersions，返回全部对象的元数据的数组
	result, err := GetAll(bucket, name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 判断长度 如果长度=0直接返回
	size := len(result)
	info := allObjInfo{int64(size), nil}
	if size == 0 {
		info.Size = 0
		b, _ := json.Marshal(info)
		w.Write(b)
		return
	}
	// 如果长度不为0 将数据分页后返回
	info = pageHelper(page, result)
	b, _ := json.Marshal(info)
	w.Write(b)
}

// GetAll 获取所有对象最新版本
func GetAll(bucket string, name string) ([]es.Metadata, error) {
	from := 0
	size := 1000
	result := make([]es.Metadata, 0)
	//无限循环
	for {
		// 通过对象名 调用es包的SearchAllVersions，返回某个对象的元数据的数组
		metas, err := es.SearchApiVersions(bucket, name, from, size)
		// 如果报错
		if err != nil {
			// 打印错误并返回500
			return result, err
		}
		result = append(result, metas...)
		// 如果长度数据长度不等于1000，此时元数据服务中没有更多的数据了
		if len(metas) != size {
			// 结束循环
			return result, nil
		}
		//否则把from的值+1000进行下一次迭代
		from += size
	}
}

var limit = 5

type allObjInfo struct {
	Size int64
	Data []es.Metadata
}

func pageHelper(page int, data []es.Metadata) allObjInfo {
	size := len(data)
	info := allObjInfo{int64(size), nil}
	if size == 0 {
		info.Size = 0
		return info
	}
	metadata := make([]es.Metadata, 0)
	//手写分页
	start := (page - 1) * limit
	end := page * limit
	if start > len(data) {
		return info
	}
	if len(data) < end {
		end = len(data)
	}
	metadata = data[start:end]
	info.Data = metadata

	return info
}
