package rs

import (
	"ceph/chapter6/lib/objectStream"
	"errors"
	"fmt"
	"io"
)

type RSPutStream struct {
	*encoder
}

//dataServers 数据服务节点的地址
//hash 需要put的对象内容的hash值
//size 内容长度
func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, errors.New("dataServers number mismatch")
	}
	//perShard:根据size计算每个分片的大小
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var err error
	for i := range writers {
		writers[i], err = objectStream.NewTempPutStream(
			dataServers[i],
			fmt.Sprintf("%s.%d", hash, i),
			perShard)
		if err != nil {
			return nil, err
		}
	}
	enc := NewEncoder(writers)
	return &RSPutStream{enc}, nil
}
func (rs *RSPutStream) Commit(success bool) {
	//首先将缓存中最后的数据都写入
	rs.Flush()
	//依次调用临时对象转正或删除
	for i := range rs.writers {
		rs.writers[i].(*objectStream.TempPutStream).Commit(success)
	}
}
