package temp

import (
	"OSS/lib/golog"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func Head(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if err != nil {
		golog.Error.Println("open file err：", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		golog.Error.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
