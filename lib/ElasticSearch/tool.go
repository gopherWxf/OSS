package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//TODO tools待更新，暂时不改，三个修复工具是不能用的
/*
	tools
*/

func HasHash(bucket, hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata_%s/_search?q=hash.subField:%s&size=0", choose(available), bucket, hash)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total.Value != 0, nil
}

//可优化限制只返回一条数据，目前弄不出来
func SearchHashSize(bucket, hash string) (size int64, err error) {
	url := fmt.Sprintf("http://%s/metadata_%s/_search?q=hash:%s&size=1",
		choose(available), bucket, hash)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search hash size: %d", resp.StatusCode)
		return 0, err
	}
	result, _ := ioutil.ReadAll(resp.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
