package rs

import (
	"ceph/chapter5/lib/objectStream"
	"fmt"
	"io"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	//检查是否满足4+2 RS码的要求
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		//如果分片id对应的数据服务节点地址为空，则说明分片丢失，随机取一个补上
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		//如果存在，则打开一个对象读取流用于读取分片数据
		reader, err := objectStream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if err == nil {
			readers[i] = reader
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var err error
	for i := range readers {
		//如果读取流为空，则创建临时对象写入流用于恢复分片
		if readers[i] == nil {
			writers[i], err = objectStream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if err != nil {
				return nil, err
			}
		}
	}
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}
func (rs *RSGetStream) Close() {
	//遍历writers成员
	for i := range rs.writers {
		//如果不为nil，说明需要转正
		if rs.writers[i] != nil {
			rs.writers[i].(*objectStream.TempPutStream).Commit(true)
		}
	}
}
