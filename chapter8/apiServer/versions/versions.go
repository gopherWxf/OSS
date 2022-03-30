package versions

/*
	http://apiServerIP/versions/ 查看所有对象的所有版本信息
	http://apiServerIP/versions/<xxx> 查看指定对象的所有版本信息
	通过es的api去构造url，返回的是json，解析到结构体，然后输出即可
	metas, err := es.SearchAllVersions(object, from, size)
	url := fmt.Sprintf("http://%s/metadata/objects/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
*/
import (
	es "ceph/chapter8/lib/ElasticSearch"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
