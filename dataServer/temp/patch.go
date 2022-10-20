package temp

/*
	将数据暂存下来，等待转正，并进行数据校验
*/
import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Patch(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("uuid", uuid)
	//获取临时对象的信息
	tempinfo, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	//打开临时对象的数据文件
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	//写入数据
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//获取数据文件的信息
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//当前文件的大小
	actual := info.Size()
	//如果当前文件的大小超过tempinfo中记录的大小
	//那么就删除数据文件和信息文件
	if actual > tempinfo.Size {
		os.Remove(dataFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func readFromFile(uuid string) (*tempInfo, error) {
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, _ := io.ReadAll(file)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}
