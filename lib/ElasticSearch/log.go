package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/* ******************************** LOG ************************************ */

// Log log结构体
type Log struct {
	OsName   string `json:"osName"`
	Level    string `json:"level"`
	DateTime int64  `json:"dateTime"`
	Content  string `json:"content"`
}

// 搜索日志的响应结构体
type searchLogResult struct {
	Hits struct {
		Total struct {
			Value    int
			Relation string
		}
		Hits []struct {
			Source Log `json:"_source"`
		}
	}
}

// 搜索日志的请求结构体
type searchLogBody struct {
	Query struct {
		Bool struct {
			Must []interface{} `json:"must"`
		} `json:"bool"`
	} `json:"query"`
	Sort []map[string]string `json:"sort"`
	From int                 `json:"from"`
	Size int                 `json:"size"`
}

// PutLog 添加日志
func PutLog(doc string) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/log/_doc",
		choose(available))

	request, _ := http.NewRequest("POST", url, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		log.Println(e)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		log.Println(fmt.Errorf("fail to put log: %d %s", r.StatusCode, string(result)))
		return
	}
}

// SearchLog 获取当天某个主机 某个级别的日志
func SearchLog(searchParam map[string]interface{}, from int, size int) ([]Log, error) {
	var logData []Log             // 结果
	var requestBody searchLogBody // 请求体
	var body = ""

	if len(searchParam) > 0 { //参数为空 查询全部
		for k, v := range searchParam {
			// 如果是查询内容则使用分词
			if k == "content" {
				s := v.(string)
				requestBody.Query.Bool.Must = append(requestBody.Query.Bool.Must, map[string]map[string]string{"match": {k: s}})
			} else if k == "dateTime" { // 如果是使用时间和日期组合查询
				// s为from：时间戳 和to：时间戳
				s := v.(map[string]interface{})
				var fromDateTime float64 = 0
				var toDataTime float64 = 0
				for sk, sv := range s {
					svString := sv.(float64)
					switch sk {
					case "from":
						fromDateTime = svString
						break
					case "to":
						toDataTime = svString
						break
					}
				}
				// 组装es请求数据
				data := map[string]map[string]map[string]float64{"range": {"dateTime": {"from": fromDateTime, "to": toDataTime}}}
				requestBody.Query.Bool.Must = append(requestBody.Query.Bool.Must, data)
			} else { // 其他字段使用强制匹配
				s := v.(string)
				requestBody.Query.Bool.Must = append(requestBody.Query.Bool.Must, map[string]map[string]string{"match_phrase": {k: s}})
			}
		}
		requestBody.Sort = append(requestBody.Sort, map[string]string{"dateTime": "desc"})
		requestBody.From = from
		requestBody.Size = size
		marshal, _ := json.Marshal(requestBody)
		body = string(marshal)
	} else {
		body = fmt.Sprintf(`
		{	
			"sort": [
				{
					"dateTime": "desc"
				}
			],
			"from":%d,
			"size": %d
		}
		`, from, size)
	}
	client := http.Client{}
	url := fmt.Sprintf("http://%s/log/_search", choose(available))

	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	r, e := client.Do(request)
	if e != nil {
		return logData, e
	}

	if r.StatusCode != http.StatusOK {
		return logData, fmt.Errorf("查询日志失败: %d", r.StatusCode)
	}
	result, _ := ioutil.ReadAll(r.Body)

	var sr searchLogResult
	json.Unmarshal(result, &sr)

	for i := range sr.Hits.Hits {
		logData = append(logData, sr.Hits.Hits[i].Source)
	}
	return logData, nil
}
