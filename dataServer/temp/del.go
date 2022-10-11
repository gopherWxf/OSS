package temp

/*
	不转正，删除两个暂存文件
*/
import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	fmt.Println(uuid, "remove")
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(dataFile)
}
