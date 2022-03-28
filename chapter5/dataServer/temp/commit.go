package temp

import (
	"ceph/chapter5/dataServer/locate"
	"ceph/chapter5/utils"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func commitTempObject(dataFile string, tempinfo *tempInfo) {
	file, _ := os.Open(dataFile)
	d := url.PathEscape(utils.CalculateHash(file))
	file.Close()
	os.Rename(dataFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name+"."+d)
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
