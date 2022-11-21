package objects

import (
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"OSS/lib/rs"
	"OSS/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Post(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucket == "" {
		golog.Error.Println("url 缺少 bucket 字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 获得对象的名字
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	//获取对象内容的大小
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取hash值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		golog.Error.Println("请求头缺少 hash 字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//如果已经存在，则直接添加元数据到ES中，并返回200ok
	if locate.Exist(url.PathEscape(hash)) {
		err = es.AddVersion(bucket, name, hash, size)
		if err != nil {
			golog.Error.Println("es add version err：", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	//如果不存在，则随机挑选6个数据节点
	dataServers, err := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if err != nil {
		golog.Error.Println("找不到足够的数据服务节点")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	//调用rs.NewRSResumablePutStream生成数据流stream
	stream, err := rs.NewRSResumablePutStream(dataServers, name, url.PathEscape(hash), size)
	if err != nil {
		golog.Error.Println("new rs stream err：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//调用ToToken生成一个字符串token，放入Location响应头中，返回201StatusCreated
	w.Header().Set("location", "/temp/"+bucket+"/"+url.PathEscape(stream.ToToken()))
	//TODO base64包含/ 而pathEscape会转化/ ,那么解析的时候怎么办呢，用	url.PathUnescape()
	w.WriteHeader(http.StatusCreated)
	golog.Info.Println("获取分片上传token成功")
}
