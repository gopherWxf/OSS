package objects

import (
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/rs"
	"OSS/utils"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	defer r.Body.Close()
	//获取对象名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionID := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionID) != 0 {
		version, err = strconv.Atoi(versionID[0])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//从es中获取对象的元数据
	meta, err := es.GetMetadata(object, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//如果哈希值为空则说明被标记为删除
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)

	size, _ := strconv.Atoi(meta.Size)

	//获取hash的读取流
	stream, err := GetStream(hash, int64(size))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//获取偏移量offset
	var offset = utils.GetOffsetFromHeader(r.Header)
	//如果不为0，那么需要调用Seek将数据流跳到offset处
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes%d-%d/%d", offset, size-1, size))
		w.WriteHeader(http.StatusPartialContent)
	}
	//如果客户端想要压缩后的数据
	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		//stream的内容写入w2，数据就会被自动压缩，压缩后的数据会被写入w
		_, err = io.Copy(w2, stream)
		w2.Close()
	} else {
		_, err = io.Copy(w, stream)
	}
	//如果发送错误，说明对象在RS解码过程中发生了错误，这对象已经无法读取了
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		stream.CloseForErr()
		return
	}
	//GET对象时会对切实的分片进行修复，如果没有发生错误，恢复分片的写入需要用到temp接口的转正，Close方法用于将写入恢复分片转正
	stream.Close()
}

//获取object文件的读取流
/*
	增加size参数是因为RS码的实现要求每一个数据片的长度完全一样，在编码时如果对象长度不能被4整除
	函数会对最后一个数据片进行填充。依次在解码时必须提供对象的准确长度，防止填充数据当初原始对象数据返回
*/
func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	//查找哪几台数据节点存了该object的数据分片
	locateInfo, err := locate.Locate(hash)
	if err != nil {
		return nil, err
	}
	//如果小于4，说明数据不完整，返回错误
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail,result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	//如果长度不等于6，说明有部分数据丢失
	//则调用ChooseRandomDataServer随机选取部分用于接收恢复分片的数据服务节点
	if len(locateInfo) != rs.ALL_SHARDS {
		//随机返回多个活跃的数据节点
		dataServers, _ = heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
