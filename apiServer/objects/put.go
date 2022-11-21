package objects

import (
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"OSS/lib/rs"
	"OSS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Put(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	//获取请求头部中内容的的hash值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		golog.Error.Println("请求头中缺少digest SHA-256字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucket == "" {
		golog.Error.Println("路径参数中缺少桶名")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//获取内容的大小
	size := utils.GetSizeFromHeader(r.Header)
	//将hash值作为数据节点存储文件的名称，实现对象名与内容的解耦
	statusCode, err := storeObject(r.Body, hash, size)
	if err != nil {
		golog.Error.Println("store object err：", err)
		w.WriteHeader(statusCode)
		return
	}
	if statusCode != http.StatusOK {
		golog.Error.Println("store object err：", err)
		w.WriteHeader(statusCode)
		return
	}
	// 获取对象名
	name := strings.Split(r.URL.EscapedPath(), "/")[3]

	//给该对象增加新版本
	err = es.AddVersion(bucket, name, hash, size)
	if err != nil {
		golog.Error.Println("es add version err：", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(statusCode)
	rdb := utils.Rds
	rdb.Incr("OssEcharts" + time.Now().Format("2006-01-02"))
	//url中文%xx 改utf-8格式
	name, _ = url.QueryUnescape(name)
	golog.Info.Println(fmt.Sprintf("上传对象 %s 成功", name))
}

//将body中的内容存入object流中  在客户端将内容传递给stream的同时，进行数据校验，如果不一致，则不转正，调用DELETE
func storeObject(r io.Reader, hash string, size int64) (int, error) {
	//是否有节点存过该hash,如果存在，则跳过后续上传，直接返回200ok
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	//随机选择多个数据节点创建一个输入流
	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	//类似于linux的tee命令，参数是io.Reader和io.Writer，返回io.Reader
	//当reader被读取时，实际内容来自于r，被读取中，同时会写入stream
	//即当utils.CalculateHash(reader)从reader中读取数据的同时也写入了stream
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	//如果计算的hash值与客户端传递的hash不一致，则删除临时对象
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	//如果一致，则将临时对象转正
	stream.Commit(true)
	return http.StatusOK, nil
}

//随机选择一个数据节点创建一个输入流
func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	//随机选多个活跃节点
	serverAddrs, err := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if err != nil {
		return nil, err
	}
	//创建一个输入流并返回
	return rs.NewRSPutStream(serverAddrs, hash, size)
}
