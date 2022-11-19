package es

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

// Metadata 元数据结构体
type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
	Time    int64
}

//获取对象的元数据信息
func getMetadata(bucket, name string, versionID int) (meta Metadata, err error) {
	//%s_%d--->name_version,根据对象的名称和版本号来获取对象的元数据
	//bucket name version -->_doc真的需要吗
	//TODO 测试删除_doc会发生什么
	url := fmt.Sprintf("http://%s/metadata_%s/_doc/%s_%d/_source",
		choose(available), bucket, name, versionID)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d : %d", name, versionID, res.StatusCode)
		return
	}
	result, _ := io.ReadAll(res.Body)
	json.Unmarshal(result, &meta)
	return
}

type hit struct {
	Source Metadata `json:"_source"`
}

//type searchResult struct {
//	Hits struct {
//		Total int   `json:"total"`
//		Hits  []hit `json:"hits"`
//	} `json:"hits"`
//}

// 搜索元数据的结构体
type searchResult struct {
	Hits struct {
		Total struct {
			Value    int
			Relation string
		}
		Hits []struct {
			Source Metadata `json:"_source"`
		}
	}
}

//获取对象最新的元数据信息
func SearchLatestVersion(bucket, name string) (meta Metadata, err error) {
	//调用es的api，指定对象的名称，以版本号降序返回 第一个结果
	url := fmt.Sprintf("http://%s/metadata_%s/_search?q=name:%s&size=1&sort=version:desc",
		choose(available), bucket, url2.PathEscape(name))

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest matadata: %d", resp.StatusCode)
		return
	}
	result, _ := io.ReadAll(resp.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

//获取对象的元数据信息，需要注意这里如果version=0，则获取对象的最新元数据信息
func GetMetadata(bucket, name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(bucket, name)
	}
	return getMetadata(bucket, name, version)
}

//插入对象的元数据信息
func PutMetadata(bucket, name string, version int, size int64, hash string) error {
	//相当于es中的一条记录
	time, _ := strconv.Atoi(time.Now().Format("20060102150405"))
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s","time":%d}`, name, version, size, hash, time)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata_%s/_doc/%s_%d?op_type=create",
		choose(available), bucket, name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	//op_type=create，多个客户端上传同一元数据，只有第一个会成功
	//别的都会发生冲突，会返回409Conflict
	//此时让版本号加1继续上传
	if res.StatusCode == http.StatusConflict {
		return PutMetadata(bucket, name, version+1, size, hash)
	}
	if res.StatusCode != http.StatusCreated {
		result, _ := io.ReadAll(res.Body)
		return fmt.Errorf("fail to put metadata: %d %s", res.StatusCode, string(result))
	}
	return nil
}

//增加一次版本号
func AddVersion(bucket, name, hash string, size int64) error {
	version, err := SearchLatestVersion(bucket, name)
	if err != nil {
		return err
	}
	return PutMetadata(bucket, name, version.Version+1, size, hash)
}

//用于获取某个对象，或者所有对象的，全部版本。from和size参数指定分页的显示结果
func SearchAllVersions(bucket, name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata_%s/_search?sort=name,version&from=%d&size=%d",
		choose(available), bucket, from, size)

	if name != "" {
		url += "&q=name:" + name
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	metas := make([]Metadata, 0)
	result, _ := io.ReadAll(resp.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

// es的聚合去重并分页
func SearchApiVersions(mapping string, name string, from int, size int) ([]Metadata, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata_%s/_search", choose(available), mapping)
	var body string
	if name != "" {
		body = fmt.Sprintf(`
		{
			"query": {
        		"match_phrase": {
					"name": "%s"
        		}
    		},
			"collapse": {
				"field": "name.subField"
			},
			"sort": [
				{
					"time":"desc"
				},
				{
					"version": "desc"
				}
			],
			"from":%d,
			"size": %d
		}
	`, name, from, size)
	} else {
		body = fmt.Sprintf(`
		{
			"query": {
				"match_all": {}
			},
			"collapse": {
				"field": "name.subField"
			},
			"sort": [
				{
					"time":"desc"
				},
				{
					"version": "desc"
				}
			],
			"from":%d,
			"size": %d
		}
	`, from, size)
	}

	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}

	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	metas := make([]Metadata, 0)
	for _, hit := range sr.Hits.Hits {
		metas = append(metas, hit.Source)
	}
	for i := range metas {
		metas[i].Name, _ = url2.QueryUnescape(metas[i].Name)
	}
	return metas, nil
}

func DelMetadata(mapping string, name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata_%s/_doc/%s_%d",
		choose(available), mapping, name, version)
	request, _ := http.NewRequest("DELETE", url, nil)

	client.Do(request)
}
