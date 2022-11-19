package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 从路径中获得节点ip
	nodeIp := strings.Split(r.URL.EscapedPath(), "/")[2]
	url := fmt.Sprintf("http://%s/systemInfo", nodeIp)
	if nodeIp == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("从 %s 获取系统信息失败", nodeIp)
		w.WriteHeader(resp.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(resp.Body)
	w.Write(result)
}
