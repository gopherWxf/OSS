package temp

import (
	"OSS/chapter10/dataServer/locate"
	"OSS/chapter10/utils"
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func commitTempObject(dataFile string, tempinfo *tempInfo) {
	file, _ := os.Open(dataFile)
	d := url.PathEscape(utils.CalculateHash(file))
	file.Seek(0, io.SeekStart)
	w, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Name + "." + d)
	defer w.Close()
	w2 := gzip.NewWriter(w)
	defer w2.Close()
	//将临时文件的数据写进gzip的writer，
	io.Copy(w2, file)
	//最后删除临时文件，添加缓存定位即可
	file.Close()
	os.Remove(dataFile)
	locate.Add(tempinfo.hash(), tempinfo.id())
}
func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}
func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}
