package es

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//TODO tools待更新，暂时不改，三个修复工具是不能用的
/*
	tools
*/

//根据name和version，调用es删除数据的api删除元数据
func DelMetaData(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	res, err := client.Do(request)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Println("status:", res.StatusCode, " err:", err)

	}
}
func HasHash(hash string) (bool, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
{
  "query": {
    "match": {
      "hash": "%s"
    }
  }
}
`, hash)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	//真nm坑爹，找了好久，fuck！！！
	request.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(request)

	b, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total != 0, nil
}

//可优化限制只返回一条数据，目前弄不出来
func SearchHashSize(hash string) (int64, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
{
  "query": {
    "match": {
      "hash": "%s"
    }
  }
}
`, hash)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(request)
	b, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	size, err := strconv.ParseInt(sr.Hits.Hits[0].Source.Size, 0, 64)
	if err != nil {
		return 0, err
	}
	return size, err
}
