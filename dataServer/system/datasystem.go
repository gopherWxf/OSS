package system

import (
	"OSS/lib/system"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	info := system.GetInfo()

	marshal, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}
