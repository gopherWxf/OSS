package system

import (
	"OSS/apiServer/versions"
	es "OSS/lib/ElasticSearch"
	RedisMQ "OSS/lib/Redis"
	"OSS/lib/golog"
	"OSS/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NodeGet(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 从路径中获得节点ip
	nodeIp := strings.Split(r.URL.EscapedPath(), "/")[2]
	url := fmt.Sprintf("http://%s/systemInfo", nodeIp)
	if nodeIp == "" {
		golog.Error.Println("system info 缺少 node ip")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		golog.Error.Println("system info err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		golog.Error.Println("获取系统信息失败")
		w.WriteHeader(resp.StatusCode)
		return
	}
	//io.Copy(w, resp.Body)
	result, _ := ioutil.ReadAll(resp.Body)
	w.Write(result)
	golog.Info.Println(fmt.Sprintf("获取 %s 节点硬件信息成功", nodeIp))
}

type Info struct {
	Obj       int64             `json:"Obj"`       //对象总数量     	遍历es即可
	Put       int64             `json:"Put"`       //上传请求次数   	累加Echarts
	Uphold    int64             `json:"Uphold"`    //维护次数    	redis string OssUpHold
	Echarts   map[string]int64  `json:"Echarts"`   //每日上传次数 	redis string OssEcharts年-月-日
	Operation RedisMQ.Operation `json:"Operation"` //历史维护信息
	// op日期--list-->op日期时间       op日期时间--string-->op
}

func UseGet(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	//给Operation使用
	index, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	system := Info{
		Obj:     getObjNum(),
		Put:     getPutNum(),
		Uphold:  upholdNum(),
		Echarts: getEcharts(),
	}
	system.Operation = *getOperation(index)
	golog.Info.Println("获取系统维护信息")
	b, _ := json.Marshal(system)
	w.Write(b)
}
func getObjNum() int64 {
	buckets := es.GetAllBucket()
	if len(buckets) == 0 {
		return 0
	}
	var ans int64
	for _, bucket := range buckets {
		metas, err := versions.GetAll(bucket, "")
		if err != nil {
			golog.Error.Println("get all bucket err：", err)
			return ans
		}
		ans += int64(len(metas))
	}
	return ans
}

func getEcharts() map[string]int64 {
	//OssEcharts日期 ---> value
	key := fmt.Sprintf("%s%d%s", "OssEcharts", time.Now().Year(), "*")
	rdb := utils.Rds
	return rdb.GetEcharts(key)
}

func getPutNum() int64 {
	info := getEcharts()
	var ans int64
	for _, v := range info {
		ans += v
	}
	return ans
}

func upholdNum() int64 {
	//OssUpHold----->val
	rdb := utils.Rds
	return rdb.GetUpHoldNum("OssUpHold")
}

func getOperation(index int) *RedisMQ.Operation {
	rdb := utils.Rds
	hash := "op"
	return rdb.GetOp(hash, index)
	//op日期--list-->op日期时间       op日期时间--string-->op
}
