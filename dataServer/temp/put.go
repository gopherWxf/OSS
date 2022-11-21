package temp

/*
	转正，将$STORAGE_ROOT/temp/t.Uuid.dat 改为 $STORAGE_ROOT/objects/hash
*/
import (
	"OSS/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func Put(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, err := readFromFile(uuid)
	if err != nil {
		golog.Error.Println("read from file err: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	file, err := os.Open(dataFile)
	if err != nil {
		golog.Error.Println("open file err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := file.Stat()
	if err != nil {
		golog.Error.Println("read file stat err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	file.Close()
	if actual != tempinfo.Size {
		os.Remove(dataFile)
		golog.Error.Println("actual size mismatch,expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//如果大小一致，则转正
	commitTempObject(dataFile, tempinfo)
}
