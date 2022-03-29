package ElasticSearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"strings"
)

type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Size    string `json:"size"`
	Hash    string `json:"hash"`
}

//获取对象的元数据信息
func getMetadata(name string, versionID int) (meta Metadata, err error) {
	//%s_%d--->name_version,根据对象的名称和版本号来获取对象的元数据
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionID)
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
type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

//获取对象最新的元数据信息
func SearchLatestVersion(name string) (meta Metadata, err error) {
	//调用es的api，指定对象的名称，以版本号降序返回 第一个结果
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url2.PathEscape(name))
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest matadata: %d", res.StatusCode)
		return
	}
	result, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

//获取对象的元数据信息，需要注意这里如果version=0，则获取对象的最新元数据信息
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

//插入对象的元数据信息
func PutMetadata(name string, version int, size int64, hash string) error {
	//相当于es中的一条记录
	doc := fmt.Sprintf(`{"name":"%s","version":"%d","size":"%d","hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
		os.Getenv("ES_SERVER"), name, version)
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
		return PutMetadata(name, version+1, size, hash)
	}
	if res.StatusCode != http.StatusCreated {
		result, _ := io.ReadAll(res.Body)
		return fmt.Errorf("fail to put metadata: %d %s", res.StatusCode, string(result))
	}
	return nil
}

//增加一次版本号
func AddVersion(name, hash string, size int64) error {
	version, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	v, _ := strconv.Atoi(version.Version)
	return PutMetadata(name, v+1, size, hash)
}

//用于获取某个对象，或者所有对象的，全部版本。from和size参数指定分页的显示结果
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/objects/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	metas := make([]Metadata, 0)
	result, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}
