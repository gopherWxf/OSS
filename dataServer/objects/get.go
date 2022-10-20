package objects

import (
	"OSS/utils"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Get(ctx *gin.Context) {
	log.Println("in data server")
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	filename := strings.Split(r.URL.EscapedPath(), "/")[2]
	file := getFile(filename)
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
}
func sendFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	//从gzipReadr读出来的数据，会被gzip先解压，然后才会被读出来
	gzipStream, err := gzip.NewReader(f)
	if err != nil {
		log.Println(err)
		return
	}
	defer gzipStream.Close()
	io.Copy(w, gzipStream)
}

//返回分片文件的名称
func getFile(name string) string {
	//查找objects目录下所有以<hash>.<X>开头的文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	//校验，如果不一致则删除数据
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Println("object hash mismatch,remove", file)
		//locate.Del(hash)
		rdb := utils.Rds
		rdb.RemoveFile(hash, os.Getenv("LISTEN_ADDRESS"))

		err := os.Remove(file)
		if err != nil {
			log.Println("os remove err:", err)
		}
		return ""
	}
	return file
}
