package locate

import (
	"OSS/lib/rs"
	"OSS/utils"
)

//判断收到的数据分片是否大于等于要求，即判断收到的反馈消息数量是否大于等于4
func Exist(hash string) bool {
	locateinfo, _ := Locate(hash)
	return len(locateinfo) >= rs.DATA_SHARDS
}

//查找哪几台数据节点存了该object的数据分片
//locateInfo map[int]string key是分片的id，val是该分片的数据节点的地址
func Locate(hash string) (locateInfo map[int]string, err error) {
	//创建一个redis连接
	rdb := utils.Rds

	locateInfo = make(map[int]string)
	result, err := rdb.GetZsetIdAndIP(hash)
	if err != nil {
		return nil, err
	}
	//将获取到的分片id和节点地址存进map中返回
	for i, v := range result {
		locateInfo[i] = v
	}
	return
}
