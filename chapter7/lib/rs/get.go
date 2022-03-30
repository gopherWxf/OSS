package rs

import (
	"ceph/chapter7/lib/objectStream"
	"fmt"
	"io"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	//检查是否满足4+2 RS码的要求，不满足则返回错误
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
		//会有一种情况导致err,readers[i]=nil , 对象读取流打开失败
		if err == nil {
			readers[i] = reader
		}
	}
	writers := make([]io.Writer, ALL_SHARDS)
	//计算每个分片的大小
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
	//处理完成后，readers和writers数组形成互补的关系，对于某个分片id来说，要么在reader读取流，要么在writer写入流
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
func (rs *RSGetStream) CloseForErr() {
	//遍历writers成员
	for i := range rs.writers {
		//如果不为nil，说明需要转正
		if rs.writers[i] != nil {
			rs.writers[i].(*objectStream.TempPutStream).Commit(false)
		}
	}
}
func (rs *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(rs, buf)
		offset -= length
	}
	return offset, nil
}
