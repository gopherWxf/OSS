package temp

/*
	不转正，删除两个暂存文件
*/
import (
	"OSS/lib/golog"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func Del(ctx *gin.Context) {
	r := ctx.Request

	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	golog.Info.Println(uuid, "remove")
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(dataFile)
}
