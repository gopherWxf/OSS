package temp

import (
	RedisMQ "OSS/lib/Redis"
	"OSS/utils"
	"compress/gzip"
	"context"
	"github.com/go-redis/redis/v8"
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

	//locate.Add(tempinfo.hash(), tempinfo.id())

	//写入redis
	rdb := RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer rdb.Client.Close()
	// ZAdd Redis `ZADD key score member [score member ...]` command.
	rdb.Client.ZAdd(context.Background(), tempinfo.hash(), &redis.Z{
		Score:  float64(tempinfo.id()),
		Member: os.Getenv("LISTEN_ADDRESS"),
	})
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
