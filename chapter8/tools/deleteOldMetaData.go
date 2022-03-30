package main

import (
	es "ceph/chapter8/lib/ElasticSearch"
	"log"
	"sort"
	"strconv"
)

/*
	删除过期元数据的工具
	同一对象，最多留存5个历史版本，最早的版本会被删掉
*/
const MIN_VERSION_COUNT = 5

func main() {
	from, size := 0, 1000
	removed := make(map[string][]int)
	for {
		//获取元数据信息
		metas, err := es.SearchAllVersions("", from, size)
		if err != nil {
			log.Println(err)
			return
		}
		for i := range metas {
			name := metas[i].Name
			version := metas[i].Version
			v, _ := strconv.Atoi(version)
			if _, ok := removed[name]; !ok {
				removed[name] = make([]int, 0)
			}
			removed[name] = append(removed[name], v)
		}
		//如果长度不等于size，说明没有更多的数据了
		if len(metas) != size {
			break
		}
		from += size
	}
	for name := range removed {
		sort.Slice(removed[name], func(i, j int) bool {
			return removed[name][i] > removed[name][j]
		})
		if len(removed[name]) <= 5 {
			continue
		}
		for _, version := range removed[name][5:] {
			es.DelMetaData(name, version)
		}
	}

}
