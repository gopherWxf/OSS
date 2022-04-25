package main

import (
	"OSS/chapter10/apiServer/objects"
	es "OSS/chapter10/lib/ElasticSearch"
	"OSS/chapter10/utils"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
	检查和修复对象数据的工具
*/
func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		realhash, _ := url.PathUnescape(hash)
		verify(realhash)
	}
}
func verify(realhash string) {
	log.Println("verify:", realhash)
	size, err := es.SearchHashSize(realhash)
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
