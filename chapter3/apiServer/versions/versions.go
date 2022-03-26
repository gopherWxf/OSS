package versions

import (
	es "ceph/chapter3/lib/ElasticSearch"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	//versions只允许GET方法
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from, size := 0, 1000
	//获取对象名
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		//获取元数据信息
		metas, err := es.SearchAllVersions(object, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			bytes, _ := json.Marshal(metas[i])
			w.Write(bytes)
			w.Write([]byte("\n"))
		}
		//如果长度不等于size，说明没有更多的数据了
		if len(metas) != size {
			return
		}
		from += size
	}
}
