package main

/*
	删除没有元数据引用的对象数据
	同一对象，最多留存5个历史版本，最早的版本会被删掉
*/
import (
	es "OSS/lib/ElasticSearch"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		realhash, _ := url.PathUnescape(hash)
		hashInMetadata, err := es.HasHash(realhash)
		if err != nil {
			log.Println(err)
			return
		}
		if !hashInMetadata {
			del(realhash)
		}
	}
}
func del(hash string) {
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	res, err := client.Do(request)
	if res.StatusCode == http.StatusOK {
		fmt.Println("delete:", hash)
	} else if err != nil {
		fmt.Println(err)
	}
}
