package temp

/*
	转正，将$STORAGE_ROOT/temp/t.Uuid.dat 改为 $STORAGE_ROOT/objects/hash
*/
import (
	"ceph/chapter4/dataServer/locate"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	file, err := os.Open(dataFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	file.Close()
	if actual != tempinfo.Size {
		os.Remove(dataFile)
		log.Println("actual size mismatch,expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//如果大小一致，则转正
	commitTempObject(dataFile, tempinfo)
}
func commitTempObject(dataFile string, tempinfo *tempInfo) {
	err := os.Rename(dataFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name)
	if err != nil {
		fmt.Println("rename err", err)
	}
	locate.Add(tempinfo.Name)
}
