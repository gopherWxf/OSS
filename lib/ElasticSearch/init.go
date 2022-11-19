package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	available []string // 可用节点
	es_server = strings.Split(os.Getenv("ES_SERVER"), ",")
)

func init() {
	available = es_server
	go connEs()
}

// es 连接响应体
type esConnResp struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUuid string `json:"cluster_uuid"`
	Version     struct {
		Number                           string    `json:"number"`
		BuildFlavor                      string    `json:"build_flavor"`
		BuildType                        string    `json:"build_type"`
		BuildHash                        string    `json:"build_hash"`
		BuildDate                        time.Time `json:"build_date"`
		BuildSnapshot                    bool      `json:"build_snapshot"`
		LuceneVersion                    string    `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

// 遍历所有节点
func connEs() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	for {
		result := make([]string, 0)
		for _, addr := range es_server {
			url := fmt.Sprintf("http://%s", addr)
			resp, err := client.Get(url)
			if err != nil {
				continue
			}
			r, _ := ioutil.ReadAll(resp.Body)
			esResp := &esConnResp{}
			err = json.Unmarshal(r, esResp)
			if err != nil {
				log.Println(err)
			}
			if esResp.Name != "" {
				result = append(result, addr)
			}
		}
		available = result
		// 延时5秒
		time.Sleep(5 * time.Second)
	}
}

// 选取一个节点
func choose(availableAddrs []string) string {
	// 从可用节点中随机选取一个节点
	rand := rand.Intn(len(available))
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	url := fmt.Sprintf("http://%s", available[rand])
	_, err := client.Get(url)
	if err != nil {
		return choose(availableAddrs)
	}
	return available[rand]
}

// AddBucket 创建映射
func AddBucket(mapping string) error {
	client := http.Client{}
	//PUT ip/metadata_bucket   索引
	url := fmt.Sprintf("http://%s/metadata_%s", choose(available), mapping)

	body := fmt.Sprintf(`
		{
    		"mappings": {
				"properties": {
            		"name": {
                		"type": "text",
                		"index": "true",
                		"fielddata":true,
                		"fields":{
                    		"subField":{
								"type":"keyword",
                        		"ignore_above":256
                    		}
                		}

            		},
            		"version": {
                		"type": "integer",
                		"index": "true"
            		},
            		"size": {
                		"type": "integer"
            		},
            		"hash": {
                		"type": "text",
                		"fields":{
                    		"subField":{
                        		"type":"keyword",
                        		"ignore_above":256
                    		}
                		}
            		}
        		}
    		}
		}
	`)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to add mapping: %d", resp.StatusCode)
		return err
	}
	return nil
}

// DelBucket 删除映射
func DelBucket(mapping string) error {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata_%s", choose(available), mapping)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to delete mapping: %d", resp.StatusCode)
		return err
	}
	return nil
}

// GetAllMapping 查找全部映射--->所有bucket
func GetAllBucket() []string {
	var buckets []string

	url := fmt.Sprintf("http://%s/_mapping", choose(available))
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("fail to mapping")
	}
	mapping, _ := ioutil.ReadAll(resp.Body)

	event := make(map[string]interface{})
	err = json.Unmarshal(mapping, &event)
	if err != nil {
		log.Println(err)
		return buckets
	}

	// 解决map无序遍历的问题
	keys := make([]string, 0, len(event))
	for key := range event {
		if key == "log" {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		//"metadata_"
		buckets = append(buckets, key[9:])
	}
	return buckets
}

// SearchBucket 搜索bucket
func SearchBucket(bucket string) int {
	url := fmt.Sprintf("http://%s/metadata_%s/_mapping", choose(available), bucket)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	return resp.StatusCode
}
