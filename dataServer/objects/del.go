package objects

import (
	"OSS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Del(ctx *gin.Context) {
	r := ctx.Request

	//获取hash
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := url.PathEscape(object)
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	//locate.Del(hash)
	rdb := utils.Rds
	defer rdb.Client.Close()
	rdb.RemoveFile(hash, os.Getenv("LISTEN_ADDRESS"))

	err := os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
	if err != nil {
		fmt.Println("rename err", err)
	}
}
