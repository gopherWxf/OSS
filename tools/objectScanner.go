package tools

import (
	"OSS/apiServer/objects"
	es "OSS/lib/ElasticSearch"
	"OSS/lib/golog"
	"OSS/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
	检查和修复对象数据的工具
*/

func verify(realhash string) {

	buckets := es.GetAllBucket()
	flag := false
	var bucket0 string
	for _, bucket := range buckets {
		hashInMetadata, err := es.HasHash(bucket, realhash)
		if err != nil {
			log.Println(err)
			return
		}
		if hashInMetadata {
			bucket0 = bucket
			flag = true
			break
		}
	}
	if flag {
		log.Println("verify:", realhash)
		size, err := es.SearchHashSize(bucket0, realhash)
		if err != nil {
			log.Println(err)
			return
		}
		//创建读取流
		stream, err := objects.GetStream(url.PathEscape(realhash), size)
		if err != nil {
			log.Println(err)
			return
		}
		d := utils.CalculateHash(stream)
		if d != realhash {
			log.Printf("object hash mismatch,calculated=%s,requested=%s\n", d, realhash)
		}
		//读取流在close的时候会把修复数据转正
		stream.Close()
	}

}
func ObjectScanner(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		realhash, _ := url.PathUnescape(hash)
		verify(realhash)
	}
	w.WriteHeader(http.StatusOK)
	rdb := utils.Rds
	rdb.Incr("OssUpHold")
	op := fmt.Sprintf("进行了对象数据全盘扫描修复的操作")
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("15:04:05")
	rdb.InsertOp(op, date, time)
	golog.Trace.Println("全盘数据扫描修复")
}
