package objects

import (
	"OSS/chapter10/dataServer/locate"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	//获取hash
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := url.PathEscape(object)
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(hash)
	err := os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
	if err != nil {
		fmt.Println("rename err", err)
	}
}
